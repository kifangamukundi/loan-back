package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/kifangamukundi/gm/libs/auths"
	"github.com/kifangamukundi/gm/libs/binders"
	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/transformations"
	"github.com/kifangamukundi/gm/loan/bindings"
	"github.com/kifangamukundi/gm/loan/emails"
	"github.com/kifangamukundi/gm/loan/models"
	"github.com/kifangamukundi/gm/loan/templates"

	"github.com/gin-gonic/gin"
)

type AgentController struct {
	AgentModel *models.AgentModel
	UserModel  *models.UserModel
}

func NewAgentController(agentModel *models.AgentModel, userModel *models.UserModel) *AgentController {
	return &AgentController{
		AgentModel: agentModel,
		UserModel:  userModel,
	}
}

func (ctrl *AgentController) CreateAgentController(c *gin.Context) {
	var req bindings.CreateAgent
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	hashedPassword, err := auths.HashPassword(req.MobileNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	user := models.User{
		FirstName:    parameters.TrimWhitespace(req.FirstName),
		LastName:     parameters.TrimWhitespace(req.LastName),
		Email:        req.Email,
		MobileNumber: req.MobileNumber,
		Password:     hashedPassword,
		IsActive:     true,
	}

	if err := ctrl.UserModel.CreateUserAgent(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	agent := models.Agent{
		UserID:    user.ID,
		IsActive:  true,
		CountryID: uint(req.CountryID),
		RegionID:  uint(req.RegionID),
		CityID:    uint(req.CityID),
	}

	if err := ctrl.AgentModel.CreateAgent(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginMode := os.Getenv("GIN_MODE")
	frontEndBaseUrl := os.Getenv("LOCAL_FRONT_END")
	if ginMode == "true" {
		frontEndBaseUrl = os.Getenv("LIVE_FRONT_END")
	}

	supportEmail := os.Getenv("SUPPORT_EMAIL")
	supportPhone := os.Getenv("SUPPORT_PHONE")
	companyName := os.Getenv("COMPANY_NAME")
	activationUrl := fmt.Sprintf("%s/login", frontEndBaseUrl)
	message := templates.GenerateAgentWelcomeMessage(user.FirstName, user.LastName, activationUrl, supportEmail, supportPhone, companyName)

	var mailProvider string
	if ginMode == "true" {
		mailProvider = "accounts"
	} else {
		mailProvider = "default"
	}

	mailOptions := map[string]string{
		"to":      user.Email,
		"subject": "Welcome to our service Agent!",
		"html":    message,
	}
	if err := emails.SendEmail(mailOptions, mailProvider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

func (ctrl *AgentController) GetAgentsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agents, totalCount, count, err := ctrl.AgentModel.GetAgents(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching agents: " + err.Error()})
		return
	}

	var transformedAgents []map[string]interface{}
	for _, agent := range agents {
		transformedAgents = append(transformedAgents, map[string]interface{}{
			"id":    agent.ID,
			"title": agent.User.FirstName + " " + agent.User.LastName,
		})
	}

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedAgents)
}

func (ctrl *AgentController) GetAllAgentsController(c *gin.Context) {
	agents, err := ctrl.AgentModel.GetAllAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting agents: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedAgents := transformations.Transform(agents, fieldNames,
		func(agent models.Agent) interface{} { return agent.ID },
		func(agent models.Agent) interface{} { return agent.User.FirstName + " " + agent.User.LastName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedAgents)
}

func (ctrl *AgentController) GetAgentByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	agent, err := ctrl.AgentModel.GetAgentByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	response := bindings.AgentResponse{
		ID:        agent.ID,
		FirstName: agent.User.FirstName,
		LastName:  agent.User.LastName,
		Email:     agent.User.Email,
		Mobile:    agent.User.MobileNumber,
		IsActive:  agent.IsActive,
	}

	binders.ReturnJSONGeneralResponse(c, response)
}

func (ctrl *AgentController) UpdateAgentController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	idInt, err := strconv.Atoi(string(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req bindings.UpdateAgent
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	agent, err := ctrl.AgentModel.UpdateAgent(idInt, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating agent: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, agent)
}

func (ctrl *AgentController) CountAgentsController(c *gin.Context) {
	conditions := map[string]interface{}{
		"is_active": true,
	}

	count, err := ctrl.AgentModel.CountAgents(conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to count active agents: %v", err)})
		return
	}

	transformedCount := map[string]interface{}{
		"count": count,
	}
	
	binders.ReturnJSONGeneralResponse(c, transformedCount)
}

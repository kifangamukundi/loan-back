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
	"github.com/kifangamukundi/gm/loan/bindings"
	"github.com/kifangamukundi/gm/loan/emails"
	"github.com/kifangamukundi/gm/loan/models"
	"github.com/kifangamukundi/gm/loan/templates"

	"github.com/gin-gonic/gin"
)

type OfficerController struct {
	OfficerModel *models.OfficerModel
	UserModel    *models.UserModel
}

func NewOfficerController(officerModel *models.OfficerModel, userModel *models.UserModel) *OfficerController {
	return &OfficerController{
		OfficerModel: officerModel,
		UserModel:    userModel,
	}
}

func (ctrl *OfficerController) CreateOfficerController(c *gin.Context) {
	var req bindings.CreateOfficer
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

	if err := ctrl.UserModel.CreateUserOfficer(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	officer := models.Officer{
		UserID:    user.ID,
		IsActive:  true,
		CountryID: uint(req.CountryID),
		RegionID:  uint(req.RegionID),
		CityID:    uint(req.CityID),
	}

	if err := ctrl.OfficerModel.CreateOfficer(&officer); err != nil {
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
	message := templates.GenerateLoanOfficerWelcomeMessage(user.FirstName, user.LastName, activationUrl, supportEmail, supportPhone, companyName)

	var mailProvider string
	if ginMode == "true" {
		mailProvider = "accounts"
	} else {
		mailProvider = "default"
	}

	mailOptions := map[string]string{
		"to":      user.Email,
		"subject": "Welcome to our service Loan Officer!",
		"html":    message,
	}
	if err := emails.SendEmail(mailOptions, mailProvider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

func (ctrl *OfficerController) GetOfficersController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	officers, totalCount, count, err := ctrl.OfficerModel.GetOfficers(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching officers: " + err.Error()})
		return
	}

	var transformedOfficers []map[string]interface{}
	for _, officer := range officers {
		transformedOfficers = append(transformedOfficers, map[string]interface{}{
			"id":    officer.ID,
			"title": officer.User.FirstName + " " + officer.User.LastName,
		})
	}

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedOfficers)
}

func (ctrl *OfficerController) GetAllOfficersController(c *gin.Context) {
	officers, err := ctrl.OfficerModel.GetAllOfficers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting officers: " + err.Error()})
		return
	}

	var transformedOfficers []map[string]interface{}
	for _, officer := range officers {
		transformedOfficers = append(transformedOfficers, map[string]interface{}{
			"id":    officer.ID,
			"title": officer.User.FirstName + " " + officer.User.LastName,
		})
	}

	binders.ReturnJSONGeneralResponse(c, transformedOfficers)
}

func (ctrl *OfficerController) GetOfficerByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	officer, err := ctrl.OfficerModel.GetOfficerByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	response := bindings.OfficerResponse{
		ID:        officer.ID,
		FirstName: officer.User.FirstName,
		LastName:  officer.User.LastName,
		Email:     officer.User.Email,
		Mobile:    officer.User.MobileNumber,
		IsActive:  officer.IsActive,
	}

	binders.ReturnJSONGeneralResponse(c, response)
}

func (ctrl *OfficerController) UpdateOfficerController(c *gin.Context) {
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

	var req bindings.UpdateOfficer
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	agent, err := ctrl.OfficerModel.UpdateOfficer(idInt, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating agent: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, agent)
}

func (ctrl *OfficerController) CountOfficersController(c *gin.Context) {
	conditions := map[string]interface{}{
		"is_active": true,
	}

	count, err := ctrl.OfficerModel.CountOfficers(conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to count active officers: %v", err)})
		return
	}

	transformedCount := map[string]interface{}{
		"count": count,
	}
	
	binders.ReturnJSONGeneralResponse(c, transformedCount)
}
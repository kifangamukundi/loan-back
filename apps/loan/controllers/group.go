package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/kifangamukundi/gm/libs/binders"
	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/transformations"
	"github.com/kifangamukundi/gm/loan/bindings"
	"github.com/kifangamukundi/gm/loan/models"

	"github.com/gin-gonic/gin"
)

type GroupController struct {
	GroupModel *models.GroupModel
	UserModel  *models.UserModel
	AgentModel *models.AgentModel
}

func NewGroupController(groupModel *models.GroupModel, userModel *models.UserModel, agentModel *models.AgentModel) *GroupController {
	return &GroupController{
		GroupModel: groupModel,
		UserModel:  userModel,
		AgentModel: agentModel,
	}
}

func (ctrl *GroupController) CreateGroupController(c *gin.Context) {
	var req bindings.CreateGroup
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	agent, err := ctrl.AgentModel.GetAgentByField("id", strconv.Itoa(req.AgentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	group := models.Group{
		GroupName: parameters.TrimWhitespace(req.GroupName),
		AgentID:   agent.ID,
		CountryID: uint(req.CountryID),
		RegionID:  uint(req.RegionID),
		CityID:    uint(req.CityID),
	}

	if err := ctrl.GroupModel.CreateGroup(&group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

func (ctrl *GroupController) GetGroupsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groups, totalCount, count, err := ctrl.GroupModel.GetGroups(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching groups: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedGroups := transformations.Transform(groups, fieldNames,
		func(group models.Group) interface{} { return group.ID },
		func(group models.Group) interface{} { return group.GroupName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedGroups)
}

func (ctrl *GroupController) GetGroupByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	group, err := ctrl.GroupModel.GetGroupByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	response := bindings.GroupResponse{
		ID:        group.ID,
		GroupName: group.GroupName,
		AgentID:   group.AgentID,
		IsActive:  group.IsActive,
	}

	binders.ReturnJSONGeneralResponse(c, response)
}

func (ctrl *GroupController) UpdateGroupController(c *gin.Context) {
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

	var req bindings.UpdateGroup
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	agent, err := ctrl.GroupModel.UpdateGroup(idInt, req.IsActive, uint(req.AgentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating agent: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, agent)
}

// testing getting groups for an agent
func (ctrl *GroupController) GetAgentGroupsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decodedUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	u := decodedUser.(models.User)

	user, err := ctrl.UserModel.GetUserByFieldPreloaded("id", strconv.Itoa(int(u.ID)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	groups, totalCount, count, err := ctrl.GroupModel.GetAgentGroups(int(user.Agent.ID), skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching groups: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"group_name",
		"created",
		"updated",
		"agent_first_name",
		"agent_last_name",
		"members",
	}

	transformGroups := transformations.Transform(groups, fieldNames,
		func(group models.Group) interface{} { return group.ID },
		func(group models.Group) interface{} { return group.GroupName },
		func(group models.Group) interface{} { return group.CreatedAt },
		func(group models.Group) interface{} { return group.UpdatedAt },
		func(group models.Group) interface{} { return group.Agent.User.FirstName },
		func(group models.Group) interface{} { return group.Agent.User.LastName },
		func(group models.Group) interface{} {
			var membersTransformed []map[string]interface{}
			for _, member := range group.Members {
				membersTransformed = append(membersTransformed, map[string]interface{}{
					"id":                member.ID,
					"member_first_name": member.User.FirstName,
					"member_last_name":  member.User.LastName,
				})
			}
			return membersTransformed
		},
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformGroups)
}

func (ctrl *GroupController) CountGroupsController(c *gin.Context) {
	conditions := map[string]interface{}{
		"is_active": true,
	}

	count, err := ctrl.GroupModel.CountGroups(conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to count groups: %v", err)})
		return
	}

	transformedCount := map[string]interface{}{
		"count": count,
	}

	binders.ReturnJSONGeneralResponse(c, transformedCount)
}

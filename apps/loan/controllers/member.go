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

type MemberController struct {
	MemberModel *models.MemberModel
	UserModel   *models.UserModel
	GroupModel  *models.GroupModel
}

func NewMemberController(memberModel *models.MemberModel, userModel *models.UserModel, groupModel *models.GroupModel) *MemberController {
	return &MemberController{
		MemberModel: memberModel,
		UserModel:   userModel,
		GroupModel:  groupModel,
	}
}

func (ctrl *MemberController) CreateMemberController(c *gin.Context) {
	var req bindings.CreateMember
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

	if err := ctrl.UserModel.CreateUserMember(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	decodedUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	u := decodedUser.(models.User)

	agent, err := ctrl.UserModel.GetUserByFieldPreloaded("id", strconv.Itoa(int(u.ID)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	_, err = ctrl.MemberModel.CreateMember(agent.Agent.ID, user.ID, true, uint(req.CountryID), uint(req.RegionID), uint(req.CityID), req.Groups)
	if err != nil {
		if err.Error() == "some groups do not exist" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating member: " + err.Error()})
		}
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
	message := templates.GenerateMemberAddedByAgentMessage(user.FirstName, user.LastName, activationUrl, supportEmail, supportPhone, companyName)

	var mailProvider string
	if ginMode == "true" {
		mailProvider = "accounts"
	} else {
		mailProvider = "default"
	}

	mailOptions := map[string]string{
		"to":      user.Email,
		"subject": "Welcome to our service Member!",
		"html":    message,
	}
	if err := emails.SendEmail(mailOptions, mailProvider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

func (ctrl *MemberController) GetGroupMembersController(c *gin.Context) {
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

	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	group, err := ctrl.GroupModel.GetGroupByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
		return
	}

	members, totalCount, count, err := ctrl.MemberModel.GetGroupMembers(int(group.ID), int(user.Agent.ID), skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching members: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"member_first_name",
		"member_last_name",
		"member_status",
		"agent_first_name",
		"agent_last_name",
		"groups",
	}

	transformedMembers := transformations.Transform(members, fieldNames,
		func(member models.Member) interface{} { return member.ID },
		func(member models.Member) interface{} { return member.User.FirstName },
		func(member models.Member) interface{} { return member.User.LastName },
		func(member models.Member) interface{} { return member.IsActive },
		func(member models.Member) interface{} { return member.Agent.User.FirstName },
		func(member models.Member) interface{} { return member.Agent.User.LastName },
		func(member models.Member) interface{} {
			var groupsTransformed []map[string]interface{}
			for _, group := range member.Groups {
				groupsTransformed = append(groupsTransformed, map[string]interface{}{
					"id":   group.ID,
					"group_name": group.GroupName,
				})
			}
			return groupsTransformed
		},
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedMembers)
}

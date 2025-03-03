package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kifangamukundi/gm/libs/auths"
	"github.com/kifangamukundi/gm/libs/binders"
	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/transformations"
	"github.com/kifangamukundi/gm/loan/bindings"
	"github.com/kifangamukundi/gm/loan/emails"
	"github.com/kifangamukundi/gm/loan/helpers"
	"github.com/kifangamukundi/gm/loan/models"
	"github.com/kifangamukundi/gm/loan/templates"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserModel *models.UserModel
}

func NewUserController(userModel *models.UserModel) *UserController {
	return &UserController{UserModel: userModel}
}

func (ctrl *UserController) RegisterUser(c *gin.Context) {
	var req bindings.RegisterRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	hashedPassword, err := auths.HashPassword(req.Password)
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
	}

	generatedToken, hashedToken, expirationDate, err := auths.GenerateAccountActivationToken(24) // 24 hours expiration
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating activation token"})
		return
	}

	user.AccountActivationToken = hashedToken
	user.AccountActivationExpire = expirationDate

	if err := ctrl.UserModel.CreateUser(&user); err != nil {
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
	activationUrl := fmt.Sprintf("%s/activate-account/%s/%d", frontEndBaseUrl, generatedToken, user.ID)
	message := templates.GenerateActivationMessage(user.FirstName, user.LastName, activationUrl, supportEmail, supportPhone, companyName)

	var mailProvider string
	if ginMode == "true" {
		mailProvider = "accounts"
	} else {
		mailProvider = "default"
	}

	mailOptions := map[string]string{
		"to":      user.Email,
		"subject": "Welcome to our service!",
		"html":    message,
	}
	if err := emails.SendEmail(mailOptions, mailProvider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Account Activation
func (ctrl *UserController) AccountActivation(c *gin.Context) {
	activationToken, validToken := parameters.ConvertParamToValidString(c, "activationToken")
	if !validToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activation token format"})
		return
	}

	id, validID := parameters.ConvertParamToValidString(c, "id")
	if !validID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := ctrl.UserModel.GetUserByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	isValid, err := auths.VerifyActivationToken(string(activationToken), user.AccountActivationToken, user.AccountActivationExpire)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired activation token"})
		return
	}

	user.IsActive = true
	user.IsLocked = false
	user.AccountActivationToken = ""
	user.AccountActivationExpire = time.Time{}

	if err := ctrl.UserModel.UpdateUserGeneric(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating user: %v", err)})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Forgot Password
func (ctrl *UserController) ForgotPassword(c *gin.Context) {
	var request bindings.ForgotPasswordRequest
	if !binders.ValidateBindJSONRequest(c, &request) {
		return
	}

	interval, err := strconv.Atoi(os.Getenv("RESET_REQUEST_INTERVAL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid interval configuration"})
		return
	}
	resetRequestInterval := time.Duration(interval) * time.Minute

	user, err := ctrl.UserModel.GetUserByField("email", request.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is inactive. Please activate your account."})
		return
	}

	if user.IsLocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is locked. Please contact support."})
		return
	}

	if user.LastResetRequestAt != (time.Time{}) {
		timeSinceLastRequest := time.Since(user.LastResetRequestAt)

		if timeSinceLastRequest < resetRequestInterval {
			timeRemaining := resetRequestInterval - timeSinceLastRequest
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Password reset request is limited to one request every %d minutes. Please wait for %d minutes.", interval, int(timeRemaining.Minutes())),
			})
			return
		}
	}

	if user.ResetRequestCount >= 3 {
		user.IsLocked = true
		if err := ctrl.UserModel.UpdateUserGeneric(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "Account has been locked due to too many password reset attempts"})
		return
	}

	generatedToken, hashedToken, expirationDate, err := auths.GenerateResetPasswordToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ResetPasswordToken = hashedToken
	user.ResetPasswordExpire = expirationDate
	user.ResetRequestCount++
	user.LastResetRequestAt = time.Now()

	if err := ctrl.UserModel.UpdateUserGeneric(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ginMode := os.Getenv("GIN_MODE")
	frontEndBaseUrl := os.Getenv("LOCAL_FRONT_END")
	if ginMode == "true" {
		frontEndBaseUrl = os.Getenv("LIVE_FRONT_END")
	}

	resetUrl := fmt.Sprintf("%s/reset-password/%s/%d", frontEndBaseUrl, generatedToken, user.ID)

	message := templates.GenerateResetPasswordMessage(user.FirstName, user.LastName, resetUrl)

	var mailProvider string
	if ginMode == "true" {
		mailProvider = "accounts"
	} else {
		mailProvider = "default"
	}

	mailOptions := map[string]string{
		"to":      user.Email,
		"subject": "Password Reset Request",
		"html":    message,
	}

	if err := emails.SendEmail(mailOptions, mailProvider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// ResetPassword handles the password reset process
func (ctrl *UserController) ResetPassword(c *gin.Context) {
	var request bindings.ChangePasswordRequest
	if !binders.ValidateBindJSONRequest(c, &request) {
		return
	}

	resetToken, validToken := parameters.ConvertParamToValidString(c, "resetToken")
	if !validToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reset token format"})
		return
	}

	id, validID := parameters.ConvertParamToValidID(c, "id")
	if !validID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := ctrl.UserModel.GetUserByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	isValid, err := auths.VerifyResetPasswordToken(string(resetToken), user.ResetPasswordToken, user.ResetPasswordExpire)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	hashedPassword, err := auths.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	user.Password = hashedPassword
	user.ResetPasswordToken = ""
	user.ResetPasswordExpire = time.Time{}

	if err := ctrl.UserModel.UpdateUserGeneric(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	companyName := os.Getenv("COMPANY_NAME")

	changedPasswordMessage := templates.GenerateChangedPasswordMessage(user.FirstName, user.LastName, companyName)

	var mailProvider string
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "true" {
		mailProvider = "accounts"
	} else {
		mailProvider = "default"
	}

	mailOptions := map[string]string{
		"to":      user.Email,
		"subject": "Your Password Has Been Changed",
		"html":    changedPasswordMessage,
	}

	if err := emails.SendEmail(mailOptions, mailProvider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset confirmation email"})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Login method
func (ctrl *UserController) Login(c *gin.Context) {
	var request bindings.LoginRequest
	if !binders.ValidateBindJSONRequest(c, &request) {
		return
	}

	email := request.Email
	password := request.Password

	user, err := ctrl.UserModel.GetUserByField("email", email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account Inactive"})
		return
	}

	if user.IsLocked {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account Locked"})
		return
	}

	isMatch, err := auths.CheckPassword(password, user.Password)
	if err != nil || !isMatch {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	userInitials := string(user.FirstName[0]) + string(user.LastName[0])

	accessToken, err := auths.GenerateAccessToken(int(user.ID), userInitials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := auths.GenerateRefreshToken(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	binders.ReturnJSONTokenResponse(c, accessToken, refreshToken)
}

// Refresh token
func (ctrl *UserController) Refresh(c *gin.Context) {
	var request bindings.RefreshTokenRequest
	if !binders.ValidateBindJSONRequest(c, &request) {
		return
	}

	claims, err := auths.VerifyRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
		return
	}

	if int(userID) != request.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token user ID mismatch"})
		return
	}

	user, err := ctrl.UserModel.GetUserByField("id", fmt.Sprintf("%d", int(userID)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userInitials := string(user.FirstName[0]) + string(user.LastName[0])

	accessToken, err := auths.GenerateAccessToken(int(user.ID), userInitials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Optionally generate a new refresh token
	// newRefreshToken, err := auths.GenerateRefreshToken(int(user.ID), user.Email)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
	// 	return
	// }

	binders.ReturnJSONTokenResponse(c, accessToken, request.RefreshToken)
}

// Addons
func (ctrl *UserController) GetUsersController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, totalCount, count, err := ctrl.UserModel.GetUsers(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedUsers := transformations.Transform(users, fieldNames,
		func(user models.User) interface{} { return user.ID },
		func(user models.User) interface{} { return user.FirstName + " " + user.LastName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedUsers)
}

func (ctrl *UserController) GetUserByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	role, err := ctrl.UserModel.GetUserByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, role)
}

func (ctrl *UserController) UpdateUserController(c *gin.Context) {
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

	var req bindings.UpdateUserRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	user, err := ctrl.UserModel.UpdateUser(idInt, req.IsActive, req.IsLocked, req.Roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, user)
}

func (ctrl *UserController) GetUserPermissionsController(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	cacheKey := helpers.GenerateUserPermissionsCacheKey("user_permissions", userModel.ID)
	ttl := 5 * time.Minute

	redisCache := helpers.InitializeRedisCache()

	cacheHit := helpers.CheckBinaryCache(redisCache, cacheKey, c, func(c *gin.Context, cachedResponse interface{}) {
		log.Printf("Data retrieved: %+v", cachedResponse)
		binders.ReturnJSONCacheResponse(c, cachedResponse)
	})

	if cacheHit {
		return
	}

	uniquePermissions := make(map[string]bool)
	for _, role := range userModel.Roles {
		for _, permission := range role.Permissions {
			uniquePermissions[permission.PermissionName] = true
		}
	}

	var permissionsList []string
	for permission := range uniquePermissions {
		permissionsList = append(permissionsList, permission)
	}

	dataFormat := map[string]interface{}{
		"permissions": permissionsList,
	}

	helpers.StoreBinaryCache(redisCache, cacheKey, dataFormat, ttl)

	binders.ReturnJSONPermissionsResponse(c, permissionsList)
}

func (ctrl *UserController) GetMonthlyUserCountsController(c *gin.Context) {
	monthlyCounts, err := ctrl.UserModel.GetMonthlyUserCounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching monthly user counts: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, monthlyCounts)
}

func (ctrl *UserController) CountUsersController(c *gin.Context) {
	conditions := map[string]interface{}{
		"is_active": true,
	}

	count, err := ctrl.UserModel.CountUsers(conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to count active users: %v", err)})
		return
	}

	transformedCount := map[string]interface{}{
		"count": count,
	}
	
	binders.ReturnJSONGeneralResponse(c, transformedCount)
}

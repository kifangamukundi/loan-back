package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, userController *controllers.UserController, db *gorm.DB) {
	forgotPasswordLimiter := rates.CreateRateLimiter("3-H")
	resetPasswordLimiter := rates.CreateRateLimiter("5-H")
	refreshLimiter := rates.CreateRateLimiter("20-H")
	updateUserLimiter := rates.CreateRateLimiter("20-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"first_name", "last_name"}
	defaultSortCriteria := "first_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/users")
	{
		v1.POST("/register", userController.RegisterUser)
		v1.PUT("/activate-account/:activationToken/:id", userController.AccountActivation)
		v1.POST("/forgot-password", forgotPasswordLimiter, userController.ForgotPassword)
		v1.PUT("/reset-password/:resetToken/:id", resetPasswordLimiter, userController.ResetPassword)
		v1.POST("/login", userController.Login)
		v1.POST("/refresh", refreshLimiter, userController.Refresh)
		v1.GET("/paginate",
			middlewares.AdvancedAuth(db, []string{"view_users"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			userController.GetUsersController,
		)
		v1.GET("/by/:id", middlewares.AdvancedAuth(db, []string{"view_users"}), userController.GetUserByIdController)
		v1.GET("/by/month", userController.GetMonthlyUserCountsController)
		v1.GET("/by/count", userController.CountUsersController)
		v1.PATCH("/by/:id", updateUserLimiter, middlewares.AdvancedAuth(db, []string{"edit_user"}), userController.UpdateUserController)
		v1.GET("/permissions", middlewares.AdvancedAuth(db, []string{}), userController.GetUserPermissionsController)
	}

	v2 := api.Group("/v2/users")
	{
		v2.POST("/register", userController.RegisterUser)
	}
}

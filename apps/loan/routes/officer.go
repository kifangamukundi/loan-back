package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OfficerRoutes(r *gin.Engine, officerController *controllers.OfficerController, db *gorm.DB) {
	createOfficerLimiter := rates.CreateRateLimiter("100-H")
	updateOfficerLimiter := rates.CreateRateLimiter("1000-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"user_id"}
	defaultSortCriteria := "user_id"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/officers")
	{
		v1.POST("/create", createOfficerLimiter, middlewares.AdvancedAuth(db, []string{"create_officer"}), officerController.CreateOfficerController)
		v1.GET("/paginate",
			middlewares.AdvancedAuth(db, []string{"view_officers"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			officerController.GetOfficersController,
		)
		v1.GET("/by/:id", middlewares.AdvancedAuth(db, []string{"view_officers"}), officerController.GetOfficerByIdController)
		v1.GET("/by/count", officerController.CountOfficersController)
		v1.PATCH("/by/:id", updateOfficerLimiter, middlewares.AdvancedAuth(db, []string{"edit_officer"}), officerController.UpdateOfficerController)
		v1.GET("/all", middlewares.AdvancedAuth(db, []string{"view_officers"}), officerController.GetAllOfficersController)
	}

	v2 := api.Group("/v2/officers")
	{
		v2.POST("/create", createOfficerLimiter, middlewares.AdvancedAuth(db, []string{"create_officer"}), officerController.CreateOfficerController)
	}
}

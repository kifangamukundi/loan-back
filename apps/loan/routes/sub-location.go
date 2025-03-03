package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SubLocationRoutes(r *gin.Engine, subLocationController *controllers.SubLocationController, db *gorm.DB) {
	createSubLocationLimiter := rates.CreateRateLimiter("100-H")
	updateSubLocationLimiter := rates.CreateRateLimiter("100-H")
	deleteSubLocationLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"sub_location_name"}
	defaultSortCriteria := "sub_location_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/sublocations")
	{
		v1.POST("/create", createSubLocationLimiter, subLocationController.CreateSubLocationController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			subLocationController.GetSubLocationsController,
		)
		v1.GET("/by/:id", subLocationController.GetSubLocationByIdController)
		v1.PATCH("/by/:id", updateSubLocationLimiter, subLocationController.UpdateSubLocationController)
		v1.DELETE("/by/:id", deleteSubLocationLimiter, subLocationController.DeleteSubLocationController)
		v1.GET("/all", subLocationController.GetAllSubLocationsController)
		v1.GET("/all/no-auth", subLocationController.GetAllSubLocationsController)
		v1.GET("/by/villages/:id", subLocationController.GetSubLocationVillagesController)
	}

	v2 := api.Group("/v2/sublocations")
	{
		v2.POST("/create", createSubLocationLimiter, subLocationController.CreateSubLocationController)
	}
}

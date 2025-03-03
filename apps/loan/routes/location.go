package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LocationRoutes(r *gin.Engine, locationController *controllers.LocationController, db *gorm.DB) {
	createLocationLimiter := rates.CreateRateLimiter("100-H")
	updateLocationLimiter := rates.CreateRateLimiter("100-H")
	deleteLocationLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"location_name"}
	defaultSortCriteria := "location_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/locations")
	{
		v1.POST("/create", createLocationLimiter, locationController.CreateLocationController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			locationController.GetLocationsController,
		)
		v1.GET("/by/:id", locationController.GetLocationByIdController)
		v1.PATCH("/by/:id", updateLocationLimiter, locationController.UpdateLocationController)
		v1.DELETE("/by/:id", deleteLocationLimiter, locationController.DeleteLocationController)
		v1.GET("/all", locationController.GetAllLocationsController)
		v1.GET("/all/no-auth", locationController.GetAllLocationsController)
		v1.GET("/by/sub-locations/:id", locationController.GetLocationSubLocationsController)
	}

	v2 := api.Group("/v2/locations")
	{
		v2.POST("/create", createLocationLimiter, locationController.CreateLocationController)
	}
}

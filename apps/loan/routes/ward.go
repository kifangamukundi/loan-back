package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WardRoutes(r *gin.Engine, wardController *controllers.WardController, db *gorm.DB) {
	createWardLimiter := rates.CreateRateLimiter("100-H")
	updateWardLimiter := rates.CreateRateLimiter("100-H")
	deleteWardLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"ward_name"}
	defaultSortCriteria := "ward_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/wards")
	{
		v1.POST("/create", createWardLimiter, wardController.CreateWardController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			wardController.GetWardsController,
		)
		v1.GET("/by/:id", wardController.GetWardByIdController)
		v1.PATCH("/by/:id", updateWardLimiter, wardController.UpdateWardController)
		v1.DELETE("/by/:id", deleteWardLimiter, wardController.DeleteWardController)
		v1.GET("/all", wardController.GetAllWardsController)
		v1.GET("/all/no-auth", wardController.GetAllWardsController)
		v1.GET("/by/locations/:id", wardController.GetWardLocationsController)
	}

	v2 := api.Group("/v2/wards")
	{
		v2.POST("/create", createWardLimiter, wardController.CreateWardController)
	}
}

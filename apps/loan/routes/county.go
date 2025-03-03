package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CountyRoutes(r *gin.Engine, countyController *controllers.CountyController, db *gorm.DB) {
	createCountyLimiter := rates.CreateRateLimiter("100-H")
	updateCountyLimiter := rates.CreateRateLimiter("100-H")
	deleteCountyLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"county_name"}
	defaultSortCriteria := "county_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/counties")
	{
		v1.POST("/create", createCountyLimiter, countyController.CreateCountyController)
		v1.GET("/paginate",

			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			countyController.GetCountysController,
		)
		v1.GET("/by/:id", countyController.GetCountyByIdController)
		v1.PATCH("/by/:id", updateCountyLimiter, countyController.UpdateCountyController)
		v1.DELETE("/by/:id", deleteCountyLimiter, countyController.DeleteCountyController)
		v1.GET("/all", countyController.GetAllCountysController)
		v1.GET("/all/no-auth", countyController.GetAllCountysController)
		v1.GET("/by/sub-counties/:id", countyController.GetCountySubCountiesController)
	}

	v2 := api.Group("/v2/counties")
	{
		v2.POST("/create", createCountyLimiter, countyController.CreateCountyController)
	}
}

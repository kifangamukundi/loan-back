package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegionRoutes(r *gin.Engine, regionController *controllers.RegionController, db *gorm.DB) {
	createRegionLimiter := rates.CreateRateLimiter("100-H")
	updateRegionLimiter := rates.CreateRateLimiter("100-H")
	deleteRegionLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"region_name"}
	defaultSortCriteria := "region_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/regions")
	{
		v1.POST("/create", regionController.CreateRegionController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			regionController.GetRegionsController,
		)
		v1.GET("/by/:id", regionController.GetRegionByIdController)
		v1.PATCH("/by/:id", updateRegionLimiter, regionController.UpdateRegionController)
		v1.DELETE("/by/:id", deleteRegionLimiter, regionController.DeleteRegionController)
		v1.GET("/all", regionController.GetAllRegionsController)
		v1.GET("/all/no-auth", regionController.GetAllRegionsController)
		v1.GET("/by/counties/:id", regionController.GetRegionCountiesController)
	}

	v2 := api.Group("/v2/regions")
	{
		v2.POST("/create", createRegionLimiter, regionController.CreateRegionController)
	}
}

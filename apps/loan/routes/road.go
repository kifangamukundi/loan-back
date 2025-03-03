package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoadRoutes(r *gin.Engine, roadController *controllers.RoadController, db *gorm.DB) {
	createRoadLimiter := rates.CreateRateLimiter("100-H")
	updateRoadLimiter := rates.CreateRateLimiter("100-H")
	deleteRoadLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"road_name"}
	defaultSortCriteria := "road_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/roads")
	{
		v1.POST("/create", createRoadLimiter, roadController.CreateRoadController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			roadController.GetRoadsController,
		)
		v1.GET("/by/:id", roadController.GetRoadByIdController)
		v1.PATCH("/by/:id", updateRoadLimiter, roadController.UpdateRoadController)
		v1.DELETE("/by/:id", deleteRoadLimiter, roadController.DeleteRoadController)
		v1.GET("/all", roadController.GetAllRoadsController)
		v1.GET("/all/no-auth", roadController.GetAllRoadsController)
		v1.GET("/by/plots/:id", roadController.GetRoadPlotsController)
	}

	v2 := api.Group("/v2/roads")
	{
		v2.POST("/create", createRoadLimiter, roadController.CreateRoadController)
	}
}

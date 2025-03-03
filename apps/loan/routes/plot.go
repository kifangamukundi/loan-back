package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PlotRoutes(r *gin.Engine, plotController *controllers.PlotController, db *gorm.DB) {
	createPlotLimiter := rates.CreateRateLimiter("100-H")
	updatePlotLimiter := rates.CreateRateLimiter("100-H")
	deletePlotLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"plot_name"}
	defaultSortCriteria := "plot_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/plots")
	{
		v1.POST("/create", createPlotLimiter, plotController.CreatePlotController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			plotController.GetPlotsController,
		)
		v1.GET("/by/:id", plotController.GetPlotByIdController)
		v1.PATCH("/by/:id", updatePlotLimiter, plotController.UpdatePlotController)
		v1.DELETE("/by/:id", deletePlotLimiter, plotController.DeletePlotController)
		v1.GET("/all", plotController.GetAllPlotsController)
		v1.GET("/all/no-auth", plotController.GetAllPlotsController)
		v1.GET("/by/units/:id", plotController.GetPlotUnitsController)
	}

	v2 := api.Group("/v2/plots")
	{
		v2.POST("/create", createPlotLimiter, plotController.CreatePlotController)
	}
}

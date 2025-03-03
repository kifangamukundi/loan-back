package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UnitRoutes(r *gin.Engine, unitController *controllers.UnitController, db *gorm.DB) {
	createUnitLimiter := rates.CreateRateLimiter("100-H")
	updateUnitLimiter := rates.CreateRateLimiter("100-H")
	deleteUnitLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"unit_name"}
	defaultSortCriteria := "unit_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/units")
	{
		v1.POST("/create", createUnitLimiter, unitController.CreateUnitController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			unitController.GetUnitsController,
		)
		v1.GET("/by/:id", unitController.GetUnitByIdController)
		v1.PATCH("/by/:id", updateUnitLimiter, unitController.UpdateUnitController)
		v1.DELETE("/by/:id", deleteUnitLimiter, unitController.DeleteUnitController)
		v1.GET("/all", unitController.GetAllUnitsController)
		v1.GET("/all/no-auth", unitController.GetAllUnitsController)
	}

	v2 := api.Group("/v2/units")
	{
		v2.POST("/create", createUnitLimiter, unitController.CreateUnitController)
	}
}

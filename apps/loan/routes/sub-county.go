package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SubCountyRoutes(r *gin.Engine, subCountyController *controllers.SubCountyController, db *gorm.DB) {
	createSubCountyLimiter := rates.CreateRateLimiter("100-H")
	updateSubCountyLimiter := rates.CreateRateLimiter("100-H")
	deleteSubCountyLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"sub_county_name"}
	defaultSortCriteria := "sub_county_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/subcounties")
	{
		v1.POST("/create", createSubCountyLimiter, subCountyController.CreateSubCountyController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			subCountyController.GetSubCountysController,
		)
		v1.GET("/by/:id", subCountyController.GetSubCountyByIdController)
		v1.PATCH("/by/:id", updateSubCountyLimiter, subCountyController.UpdateSubCountyController)
		v1.DELETE("/by/:id", deleteSubCountyLimiter, subCountyController.DeleteSubCountyController)
		v1.GET("/all", subCountyController.GetAllSubCountysController)
		v1.GET("/all/no-auth", subCountyController.GetAllSubCountysController)
		v1.GET("/by/wards/:id", subCountyController.GetSubCountyWardsController)
	}

	v2 := api.Group("/v2/subCountys")
	{
		v2.POST("/create", createSubCountyLimiter, subCountyController.CreateSubCountyController)
	}
}

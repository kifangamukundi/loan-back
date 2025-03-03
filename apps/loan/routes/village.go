package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VillageRoutes(r *gin.Engine, villageController *controllers.VillageController, db *gorm.DB) {
	createVillageLimiter := rates.CreateRateLimiter("100-H")
	updateVillageLimiter := rates.CreateRateLimiter("100-H")
	deleteVillageLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"village_name"}
	defaultSortCriteria := "village_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/villages")
	{
		v1.POST("/create", createVillageLimiter, villageController.CreateVillageController)
		v1.GET("/paginate",
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			villageController.GetVillagesController,
		)
		v1.GET("/by/:id", villageController.GetVillageByIdController)
		v1.PATCH("/by/:id", updateVillageLimiter, villageController.UpdateVillageController)
		v1.DELETE("/by/:id", deleteVillageLimiter, villageController.DeleteVillageController)
		v1.GET("/all", villageController.GetAllVillagesController)
		v1.GET("/all/no-auth", villageController.GetAllVillagesController)
		v1.GET("/by/roads/:id", villageController.GetVillageRoadsController)
	}

	v2 := api.Group("/v2/villages")
	{
		v2.POST("/create", createVillageLimiter, villageController.CreateVillageController)
	}
}

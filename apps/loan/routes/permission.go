package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PermissionRoutes(r *gin.Engine, permissionController *controllers.PermissionController, db *gorm.DB) {
	createPermissionLimiter := rates.CreateRateLimiter("100-H")
	updatePermissionLimiter := rates.CreateRateLimiter("100-H")
	deletePermissionLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"permission_name"}
	defaultSortCriteria := "permission_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/permissions")
	{
		v1.POST("/create", createPermissionLimiter, middlewares.AdvancedAuth(db, []string{"create_permission"}), permissionController.CreatePermissionController)
		v1.GET("/paginate",
			middlewares.AdvancedAuth(db, []string{"view_permissions"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			permissionController.GetPermissionsController,
		)
		v1.GET("/by/:id", middlewares.AdvancedAuth(db, []string{"view_permissions"}), permissionController.GetPermissionByIdController)
		v1.GET("/by/roles", permissionController.GetAllPermissionsByPermissionsController)
		v1.GET("/by/count", permissionController.CountPermissionsController)
		v1.PATCH("/by/:id", updatePermissionLimiter, middlewares.AdvancedAuth(db, []string{"edit_permission"}), permissionController.UpdatePermissionController)
		v1.DELETE("/by/:id", deletePermissionLimiter, middlewares.AdvancedAuth(db, []string{"delete_permission"}), permissionController.DeletePermissionController)
		v1.GET("/all", middlewares.AdvancedAuth(db, []string{"view_permissions"}), permissionController.GetAllPermissionsController)
	}

	v2 := api.Group("/v2/permissions")
	{
		v2.POST("/create", createPermissionLimiter, middlewares.AdvancedAuth(db, []string{"create_permission"}), permissionController.CreatePermissionController)
	}
}

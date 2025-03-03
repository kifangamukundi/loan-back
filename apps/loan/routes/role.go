package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoleRoutes(r *gin.Engine, roleController *controllers.RoleController, db *gorm.DB) {
	createRoleLimiter := rates.CreateRateLimiter("100-H")
	updateRoleLimiter := rates.CreateRateLimiter("100-H")
	deleteRoleLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"role_name"}
	defaultSortCriteria := "role_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/roles")
	{
		v1.POST("/create", createRoleLimiter, middlewares.AdvancedAuth(db, []string{"create_role"}), roleController.CreateRoleController)
		v1.GET("/paginate",
			middlewares.AdvancedAuth(db, []string{"view_roles"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			roleController.GetRolesController,
		)
		v1.GET("/by/:id", middlewares.AdvancedAuth(db, []string{"view_roles"}), roleController.GetRoleByIdController)
		v1.GET("/by/users", roleController.GetAllUsersByRoleController)
		v1.GET("/by/count", roleController.CountRolesController)
		v1.PATCH("/by/:id", updateRoleLimiter, middlewares.AdvancedAuth(db, []string{"edit_role"}), roleController.UpdateRoleController)
		v1.DELETE("/by/:id", deleteRoleLimiter, middlewares.AdvancedAuth(db, []string{"delete_role"}), roleController.DeleteRoleController)
		v1.GET("/all", middlewares.AdvancedAuth(db, []string{"view_roles"}), roleController.GetAllRolesController)
	}

	v2 := api.Group("/v2/roles")
	{
		v2.POST("/roles", createRoleLimiter, middlewares.AdvancedAuth(db, []string{"create_role"}), roleController.CreateRoleController)
	}
}

package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GroupRoutes(r *gin.Engine, groupController *controllers.GroupController, db *gorm.DB) {
	createGroupLimiter := rates.CreateRateLimiter("100-H")
	updateGroupLimiter := rates.CreateRateLimiter("1000-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"group_name"}
	defaultSortCriteria := "group_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/groups")
	{
		v1.POST("/create", createGroupLimiter, middlewares.AdvancedAuth(db, []string{"create_group"}), groupController.CreateGroupController)
		v1.GET("/paginate",
			middlewares.AdvancedAuth(db, []string{"view_groups"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			groupController.GetGroupsController,
		)
		v1.GET("/paginate/my",
			middlewares.AdvancedAuth(db, []string{}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			groupController.GetAgentGroupsController,
		)
		v1.GET("/by/count", groupController.CountGroupsController)
		v1.GET("/by/:id", middlewares.AdvancedAuth(db, []string{"view_groups"}), groupController.GetGroupByIdController)
		v1.PATCH("/by/:id", updateGroupLimiter, middlewares.AdvancedAuth(db, []string{"edit_group"}), groupController.UpdateGroupController)
	}

	v2 := api.Group("/v2/groups")
	{
		v2.POST("/create", createGroupLimiter, middlewares.AdvancedAuth(db, []string{"create_group"}), groupController.CreateGroupController)
	}
}

package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MemberRoutes(r *gin.Engine, memberController *controllers.MemberController, db *gorm.DB) {
	createMemberLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"user_id"}
	defaultSortCriteria := "user_id"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/members")
	{
		v1.POST("/create", createMemberLimiter, middlewares.AdvancedAuth(db, []string{"create_member"}), memberController.CreateMemberController)
		v1.GET("/paginate/:id",
			middlewares.AdvancedAuth(db, []string{"view_members"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			memberController.GetGroupMembersController,
		)
	}

	v2 := api.Group("/v2/members")
	{
		v2.POST("/create", createMemberLimiter, middlewares.AdvancedAuth(db, []string{"create_member"}), memberController.CreateMemberController)
	}
}

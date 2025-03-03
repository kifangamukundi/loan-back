package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AgentRoutes(r *gin.Engine, agentController *controllers.AgentController, db *gorm.DB) {
	createAgentLimiter := rates.CreateRateLimiter("100-H")
	updateAgentLimiter := rates.CreateRateLimiter("1000-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"user_id"}
	defaultSortCriteria := "user_id"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/agents")
	{
		v1.POST("/create", createAgentLimiter, middlewares.AdvancedAuth(db, []string{"create_agent"}), agentController.CreateAgentController)
		v1.GET("/paginate",
			middlewares.AdvancedAuth(db, []string{"view_agents"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			agentController.GetAgentsController,
		)
		v1.GET("/by/count", agentController.CountAgentsController)
		v1.GET("/by/:id", middlewares.AdvancedAuth(db, []string{"view_agents"}), agentController.GetAgentByIdController)
		v1.PATCH("/by/:id", updateAgentLimiter, middlewares.AdvancedAuth(db, []string{"edit_agent"}), agentController.UpdateAgentController)
		v1.GET("/all", middlewares.AdvancedAuth(db, []string{"view_agents"}), agentController.GetAllAgentsController)
	}

	v2 := api.Group("/v2/agents")
	{
		v2.POST("/create", createAgentLimiter, middlewares.AdvancedAuth(db, []string{"create_agent"}), agentController.CreateAgentController)
	}
}

package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoanRoutes(r *gin.Engine, loanController *controllers.LoanController, db *gorm.DB) {
	createLoanLimiter := rates.CreateRateLimiter("100-H")
	approveLoanLimiter := rates.CreateRateLimiter("100-H")
	rejectLoanLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"status"}
	defaultSortCriteria := "status"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/loans")
	{
		v1.POST("/create", createLoanLimiter, middlewares.AdvancedAuth(db, []string{"create_loan"}), loanController.CreateLoanController)
		v1.GET("/paginate/:groupId/:memberId",
			middlewares.AdvancedAuth(db, []string{"view_loans"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			loanController.GetAgentGroupMemberLoansController,
		)
		v1.GET("/paginate",
			middlewares.AdvancedAuth(db, []string{"view_loans"}),
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			loanController.GetLoansController,
		)
		v1.GET("/by/:id", middlewares.AdvancedAuth(db, []string{"view_loans"}), loanController.GetLoanByIdController)
		v1.PATCH("/by/approve/:id", approveLoanLimiter, middlewares.AdvancedAuth(db, []string{"edit_loan"}), loanController.ApproveLoanController)
		v1.PATCH("/by/reject/:id", rejectLoanLimiter, middlewares.AdvancedAuth(db, []string{"edit_loan"}), loanController.RejectLoanController)
	}

	v2 := api.Group("/v2/loans")
	{
		v2.POST("/create", createLoanLimiter, middlewares.AdvancedAuth(db, []string{"create_loan"}), loanController.CreateLoanController)
	}
}

package routes

import (
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CountryRoutes(r *gin.Engine, countryController *controllers.CountryController, db *gorm.DB) {
	createCountryLimiter := rates.CreateRateLimiter("100-H")
	createCountriesLimiter := rates.CreateRateLimiter("100-H")
	updateCountryLimiter := rates.CreateRateLimiter("100-H")
	deleteCountryLimiter := rates.CreateRateLimiter("100-H")

	validSortOrders := []string{"asc", "desc"}
	validSortCriteria := []string{"country_name"}
	defaultSortCriteria := "country_name"
	defaultPage := 1
	defaultLimit := 9

	api := r.Group("/api")

	v1 := api.Group("/v1/countries")
	{
		v1.POST("/create", createCountryLimiter, countryController.CreateCountryController)
		v1.GET("/paginate",
			createCountriesLimiter,
			queryparams.SortOrderMiddleware(validSortOrders),
			queryparams.SortColumnMiddleware(validSortCriteria, defaultSortCriteria),
			queryparams.PaginationMiddleware(defaultPage, defaultLimit),
			queryparams.SearchMiddleware(),
			queryparams.FilterMiddleware(),
			countryController.GetCountriesController,
		)
		v1.GET("/by/:id", countryController.GetCountryByIdController)
		v1.PATCH("/by/:id", updateCountryLimiter, countryController.UpdateCountryController)
		v1.DELETE("/by/:id", deleteCountryLimiter, countryController.DeleteCountryController)
		v1.GET("/all", countryController.GetAllCountriesController)
		v1.GET("/all/no-auth", countryController.GetAllCountriesController)
		v1.GET("/by/regions/:id", countryController.GetCountryRegionsController)
	}

	v2 := api.Group("/v2/countries")
	{
		v2.POST("/create", createCountryLimiter, countryController.CreateCountryController)
	}
}

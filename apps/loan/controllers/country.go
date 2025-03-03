package controllers

import (
	"net/http"
	"strconv"

	"github.com/kifangamukundi/gm/libs/binders"
	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/transformations"
	"github.com/kifangamukundi/gm/loan/bindings"
	"github.com/kifangamukundi/gm/loan/models"

	"github.com/gin-gonic/gin"
)

type CountryController struct {
	CountryModel *models.CountryModel
}

func NewCountryController(countryModel *models.CountryModel) *CountryController {
	return &CountryController{CountryModel: countryModel}
}

// Creates a new Country with body {countryName}
func (ctrl *CountryController) CreateCountryController(c *gin.Context) {
	var req bindings.CreateCountryRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	country := models.Country{
		CountryName: req.CountryName,
	}

	if err := ctrl.CountryModel.CreateCountry(&country); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of countries with paginations
func (ctrl *CountryController) GetCountriesController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	countries, totalCount, count, err := ctrl.CountryModel.GetCountries(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching countries: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedCountries := transformations.Transform(countries, fieldNames,
		func(country models.Country) interface{} { return country.ID },
		func(country models.Country) interface{} { return country.CountryName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedCountries)
}

// Get country with param {id}
func (ctrl *CountryController) GetCountryByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	country, err := ctrl.CountryModel.GetCountryByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, country)
}

// Update country with body {countryName} and param {id}
func (ctrl *CountryController) UpdateCountryController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	idInt, err := strconv.Atoi(string(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req bindings.UpdateCountryRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	country, err := ctrl.CountryModel.UpdateCountry(idInt, req.CountryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating country: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, country)
}

// Delete country with param {id}
func (ctrl *CountryController) DeleteCountryController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	idUint, err := strconv.ParseUint(string(id), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = ctrl.CountryModel.DeleteCountry(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting country: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all countries formatted in {id, title}
func (ctrl *CountryController) GetAllCountriesController(c *gin.Context) {
	countries, err := ctrl.CountryModel.GetAllCountries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting countries: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedCountries := transformations.Transform(countries, fieldNames,
		func(country models.Country) interface{} { return country.ID },
		func(country models.Country) interface{} { return country.CountryName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedCountries)
}

// Fix this so that we formatted regions though we are fetching the country itself with regions
func (ctrl *CountryController) GetCountryRegionsController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	country, err := ctrl.CountryModel.GetCountryRegionsByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedRegions := transformations.Transform(country.Regions, fieldNames,
		func(region models.Region) interface{} { return region.ID },
		func(region models.Region) interface{} { return region.RegionName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedRegions)
}

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

type CountyController struct {
	CountyModel *models.CountyModel
}

func NewCountyController(countyModel *models.CountyModel) *CountyController {
	return &CountyController{CountyModel: countyModel}
}

// Create a county with body {countyName, country_id}
func (ctrl *CountyController) CreateCountyController(c *gin.Context) {
	var req bindings.CreateCountyRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	county := models.County{
		CountyName: parameters.TrimWhitespace(req.CountyName),
		RegionID:  uint(req.RegionID),
	}

	if err := ctrl.CountyModel.CreateCounty(&county); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of countys with pagination
func (ctrl *CountyController) GetCountysController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	countys, totalCount, count, err := ctrl.CountyModel.GetCountys(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching countys: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedCountys := transformations.Transform(countys, fieldNames,
		func(county models.County) interface{} { return county.ID },
		func(county models.County) interface{} { return county.CountyName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedCountys)
}

// Get a county by param {id}
func (ctrl *CountyController) GetCountyByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	county, err := ctrl.CountyModel.GetCountyByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "County not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, county)
}

// Update a county with body {countyName, country_id} and param {id}
func (ctrl *CountyController) UpdateCountyController(c *gin.Context) {
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

	var req bindings.UpdateCountyRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.CountyName)

	county, err := ctrl.CountyModel.UpdateCounty(idInt, name, uint(req.RegionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating county: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, county)
}

// Delete county by param {id}
func (ctrl *CountyController) DeleteCountyController(c *gin.Context) {
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

	err = ctrl.CountyModel.DeleteCounty(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting county: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all countys formatted in {id and title}
func (ctrl *CountyController) GetAllCountysController(c *gin.Context) {
	countys, err := ctrl.CountyModel.GetAllCountys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting countys: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedCountys := transformations.Transform(countys, fieldNames,
		func(county models.County) interface{} { return county.ID },
		func(county models.County) interface{} { return county.CountyName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedCountys)
}

func (ctrl *CountyController) GetCountySubCountiesController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	county, err := ctrl.CountyModel.GetCountySubCountiesByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "County not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedSubCounties := transformations.Transform(county.SubCounties, fieldNames,
		func(county models.SubCounty) interface{} { return county.ID },
		func(county models.SubCounty) interface{} { return county.SubCountyName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedSubCounties)
}

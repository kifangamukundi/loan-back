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

type RegionController struct {
	RegionModel *models.RegionModel
}

func NewRegionController(regionModel *models.RegionModel) *RegionController {
	return &RegionController{RegionModel: regionModel}
}

// Create a region with body {regionName, country_id}
func (ctrl *RegionController) CreateRegionController(c *gin.Context) {
	var req bindings.CreateRegionRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	region := models.Region{
		RegionName: parameters.TrimWhitespace(req.RegionName),
		CountryID:  uint(req.CountryID),
	}

	if err := ctrl.RegionModel.CreateRegion(&region); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of regions with pagination
func (ctrl *RegionController) GetRegionsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	regions, totalCount, count, err := ctrl.RegionModel.GetRegions(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching regions: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedRegions := transformations.Transform(regions, fieldNames,
		func(region models.Region) interface{} { return region.ID },
		func(region models.Region) interface{} { return region.RegionName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedRegions)
}

// Get a region by param {id}
func (ctrl *RegionController) GetRegionByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	region, err := ctrl.RegionModel.GetRegionByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, region)
}

// Update a region with body {regionName, country_id} and param {id}
func (ctrl *RegionController) UpdateRegionController(c *gin.Context) {
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

	var req bindings.UpdateRegionRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.RegionName)

	region, err := ctrl.RegionModel.UpdateRegion(idInt, name, uint(req.CountryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating region: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, region)
}

// Delete region by param {id}
func (ctrl *RegionController) DeleteRegionController(c *gin.Context) {
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

	err = ctrl.RegionModel.DeleteRegion(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting region: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all regions formatted in {id and title}
func (ctrl *RegionController) GetAllRegionsController(c *gin.Context) {
	regions, err := ctrl.RegionModel.GetAllRegions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting regions: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedRegions := transformations.Transform(regions, fieldNames,
		func(region models.Region) interface{} { return region.ID },
		func(region models.Region) interface{} { return region.RegionName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedRegions)
}

func (ctrl *RegionController) GetRegionCountiesController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	region, err := ctrl.RegionModel.GetRegionCountiesByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Region not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedCounties := transformations.Transform(region.Counties, fieldNames,
		func(county models.County) interface{} { return county.ID },
		func(county models.County) interface{} { return county.CountyName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedCounties)
}

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

type SubLocationController struct {
	SubLocationModel *models.SubLocationModel
}

func NewSubLocationController(subLocationModel *models.SubLocationModel) *SubLocationController {
	return &SubLocationController{SubLocationModel: subLocationModel}
}

// Create a subLocation with body {subLocationName, country_id}
func (ctrl *SubLocationController) CreateSubLocationController(c *gin.Context) {
	var req bindings.CreateSubLocationRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	subLocation := models.SubLocation{
		SubLocationName: parameters.TrimWhitespace(req.SubLocationName),
		LocationID:       uint(req.LocationID),
	}

	if err := ctrl.SubLocationModel.CreateSubLocation(&subLocation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of subLocations with pagination
func (ctrl *SubLocationController) GetSubLocationsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subLocations, totalCount, count, err := ctrl.SubLocationModel.GetSubLocations(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching subLocations: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedSubLocations := transformations.Transform(subLocations, fieldNames,
		func(subLocation models.SubLocation) interface{} { return subLocation.ID },
		func(subLocation models.SubLocation) interface{} { return subLocation.SubLocationName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedSubLocations)
}

// Get a subLocation by param {id}
func (ctrl *SubLocationController) GetSubLocationByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	subLocation, err := ctrl.SubLocationModel.GetSubLocationByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubLocation not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, subLocation)
}

// Update a subLocation with body {subLocationName, country_id} and param {id}
func (ctrl *SubLocationController) UpdateSubLocationController(c *gin.Context) {
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

	var req bindings.UpdateSubLocationRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.SubLocationName)

	subLocation, err := ctrl.SubLocationModel.UpdateSubLocation(idInt, name, uint(req.LocationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating subLocation: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, subLocation)
}

// Delete subLocation by param {id}
func (ctrl *SubLocationController) DeleteSubLocationController(c *gin.Context) {
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

	err = ctrl.SubLocationModel.DeleteSubLocation(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting subLocation: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all subLocations formatted in {id and title}
func (ctrl *SubLocationController) GetAllSubLocationsController(c *gin.Context) {
	subLocations, err := ctrl.SubLocationModel.GetAllSubLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting subLocations: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedSubLocations := transformations.Transform(subLocations, fieldNames,
		func(subLocation models.SubLocation) interface{} { return subLocation.ID },
		func(subLocation models.SubLocation) interface{} { return subLocation.SubLocationName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedSubLocations)
}

func (ctrl *SubLocationController) GetSubLocationVillagesController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	subLocation, err := ctrl.SubLocationModel.GetSubLocationVillagesByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubLocation not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedVillages := transformations.Transform(subLocation.Villages, fieldNames,
		func(village models.Village) interface{} { return village.ID },
		func(village models.Village) interface{} { return village.VillageName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedVillages)
}

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

type LocationController struct {
	LocationModel *models.LocationModel
}

func NewLocationController(locationModel *models.LocationModel) *LocationController {
	return &LocationController{LocationModel: locationModel}
}

// Create a location with body {locationName, country_id}
func (ctrl *LocationController) CreateLocationController(c *gin.Context) {
	var req bindings.CreateLocationRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	location := models.Location{
		LocationName: parameters.TrimWhitespace(req.LocationName),
		WardID:  uint(req.WardID),
	}

	if err := ctrl.LocationModel.CreateLocation(&location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of locations with pagination
func (ctrl *LocationController) GetLocationsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	locations, totalCount, count, err := ctrl.LocationModel.GetLocations(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching locations: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedLocations := transformations.Transform(locations, fieldNames,
		func(location models.Location) interface{} { return location.ID },
		func(location models.Location) interface{} { return location.LocationName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedLocations)
}

// Get a location by param {id}
func (ctrl *LocationController) GetLocationByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	location, err := ctrl.LocationModel.GetLocationByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, location)
}

// Update a location with body {locationName, country_id} and param {id}
func (ctrl *LocationController) UpdateLocationController(c *gin.Context) {
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

	var req bindings.UpdateLocationRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.LocationName)

	location, err := ctrl.LocationModel.UpdateLocation(idInt, name, uint(req.WardID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating location: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, location)
}

// Delete location by param {id}
func (ctrl *LocationController) DeleteLocationController(c *gin.Context) {
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

	err = ctrl.LocationModel.DeleteLocation(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting location: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all locations formatted in {id and title}
func (ctrl *LocationController) GetAllLocationsController(c *gin.Context) {
	locations, err := ctrl.LocationModel.GetAllLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting locations: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedLocations := transformations.Transform(locations, fieldNames,
		func(location models.Location) interface{} { return location.ID },
		func(location models.Location) interface{} { return location.LocationName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedLocations)
}

func (ctrl *LocationController) GetLocationSubLocationsController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	location, err := ctrl.LocationModel.GetLocationSubLocationsByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedSubLocations := transformations.Transform(location.SubLocations, fieldNames,
		func(subLocation models.SubLocation) interface{} { return subLocation.ID },
		func(subLocation models.SubLocation) interface{} { return subLocation.SubLocationName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedSubLocations)
}

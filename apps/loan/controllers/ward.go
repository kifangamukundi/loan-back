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

type WardController struct {
	WardModel *models.WardModel
}

func NewWardController(wardModel *models.WardModel) *WardController {
	return &WardController{WardModel: wardModel}
}

// Create a ward with body {wardName, country_id}
func (ctrl *WardController) CreateWardController(c *gin.Context) {
	var req bindings.CreateWardRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	ward := models.Ward{
		WardName: parameters.TrimWhitespace(req.WardName),
		SubCountyID:  uint(req.SubCountyID),
	}

	if err := ctrl.WardModel.CreateWard(&ward); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of wards with pagination
func (ctrl *WardController) GetWardsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wards, totalCount, count, err := ctrl.WardModel.GetWards(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching wards: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedWards := transformations.Transform(wards, fieldNames,
		func(ward models.Ward) interface{} { return ward.ID },
		func(ward models.Ward) interface{} { return ward.WardName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedWards)
}

// Get a ward by param {id}
func (ctrl *WardController) GetWardByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ward, err := ctrl.WardModel.GetWardByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ward not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, ward)
}

// Update a ward with body {wardName, country_id} and param {id}
func (ctrl *WardController) UpdateWardController(c *gin.Context) {
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

	var req bindings.UpdateWardRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.WardName)

	ward, err := ctrl.WardModel.UpdateWard(idInt, name, uint(req.SubCountyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating ward: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, ward)
}

// Delete ward by param {id}
func (ctrl *WardController) DeleteWardController(c *gin.Context) {
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

	err = ctrl.WardModel.DeleteWard(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting ward: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all wards formatted in {id and title}
func (ctrl *WardController) GetAllWardsController(c *gin.Context) {
	wards, err := ctrl.WardModel.GetAllWards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting wards: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedWards := transformations.Transform(wards, fieldNames,
		func(ward models.Ward) interface{} { return ward.ID },
		func(ward models.Ward) interface{} { return ward.WardName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedWards)
}

func (ctrl *WardController) GetWardLocationsController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ward, err := ctrl.WardModel.GetWardLocationsByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ward not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedLocations := transformations.Transform(ward.Locations, fieldNames,
		func(location models.Location) interface{} { return location.ID },
		func(location models.Location) interface{} { return location.LocationName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedLocations)
}

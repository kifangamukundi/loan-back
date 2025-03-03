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

type UnitController struct {
	UnitModel *models.UnitModel
}

func NewUnitController(unitModel *models.UnitModel) *UnitController {
	return &UnitController{UnitModel: unitModel}
}

// Create unit with body { unitName and region_id}
func (ctrl *UnitController) CreateUnitController(c *gin.Context) {
	var req bindings.CreateUnitRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	unit := models.Unit{
		UnitName: req.UnitName,
		PlotID: uint(req.PlotID),
	}

	if err := ctrl.UnitModel.CreateUnit(&unit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of cities with paginations
func (ctrl *UnitController) GetUnitsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cities, totalCount, count, err := ctrl.UnitModel.GetUnits(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching cities: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedUnits := transformations.Transform(cities, fieldNames,
		func(unit models.Unit) interface{} { return unit.ID },
		func(unit models.Unit) interface{} { return unit.UnitName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedUnits)
}

// Get unit by param {id}
func (ctrl *UnitController) GetUnitByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	unit, err := ctrl.UnitModel.GetUnitByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, unit)
}

// Update unit with body {unitName and region_id} and param {id}
func (ctrl *UnitController) UpdateUnitController(c *gin.Context) {
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

	var req bindings.UpdateUnitRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	unit, err := ctrl.UnitModel.UpdateUnit(idInt, req.UnitName, uint(req.PlotID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating unit: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, unit)
}

// Delete unit with param {id}
func (ctrl *UnitController) DeleteUnitController(c *gin.Context) {
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

	err = ctrl.UnitModel.DeleteUnit(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting unit: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all cities formatted
func (ctrl *UnitController) GetAllUnitsController(c *gin.Context) {
	cities, err := ctrl.UnitModel.GetAllUnits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting cities: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedUnits := transformations.Transform(cities, fieldNames,
		func(unit models.Unit) interface{} { return unit.ID },
		func(unit models.Unit) interface{} { return unit.UnitName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedUnits)
}

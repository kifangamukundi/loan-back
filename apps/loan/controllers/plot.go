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

type PlotController struct {
	PlotModel *models.PlotModel
}

func NewPlotController(plotModel *models.PlotModel) *PlotController {
	return &PlotController{PlotModel: plotModel}
}

// Create a plot with body {plotName, country_id}
func (ctrl *PlotController) CreatePlotController(c *gin.Context) {
	var req bindings.CreatePlotRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	plot := models.Plot{
		PlotName: parameters.TrimWhitespace(req.PlotName),
		RoadID:  uint(req.RoadID),
	}

	if err := ctrl.PlotModel.CreatePlot(&plot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of plots with pagination
func (ctrl *PlotController) GetPlotsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plots, totalCount, count, err := ctrl.PlotModel.GetPlots(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching plots: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedPlots := transformations.Transform(plots, fieldNames,
		func(plot models.Plot) interface{} { return plot.ID },
		func(plot models.Plot) interface{} { return plot.PlotName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedPlots)
}

// Get a plot by param {id}
func (ctrl *PlotController) GetPlotByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	plot, err := ctrl.PlotModel.GetPlotByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plot not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, plot)
}

// Update a plot with body {plotName, country_id} and param {id}
func (ctrl *PlotController) UpdatePlotController(c *gin.Context) {
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

	var req bindings.UpdatePlotRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.PlotName)

	plot, err := ctrl.PlotModel.UpdatePlot(idInt, name, uint(req.RoadID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating plot: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, plot)
}

// Delete plot by param {id}
func (ctrl *PlotController) DeletePlotController(c *gin.Context) {
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

	err = ctrl.PlotModel.DeletePlot(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting plot: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all plots formatted in {id and title}
func (ctrl *PlotController) GetAllPlotsController(c *gin.Context) {
	plots, err := ctrl.PlotModel.GetAllPlots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting plots: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedPlots := transformations.Transform(plots, fieldNames,
		func(plot models.Plot) interface{} { return plot.ID },
		func(plot models.Plot) interface{} { return plot.PlotName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedPlots)
}

func (ctrl *PlotController) GetPlotUnitsController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	plot, err := ctrl.PlotModel.GetPlotUnitsByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plot not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedUnits := transformations.Transform(plot.Units, fieldNames,
		func(unit models.Unit) interface{} { return unit.ID },
		func(unit models.Unit) interface{} { return unit.UnitName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedUnits)
}

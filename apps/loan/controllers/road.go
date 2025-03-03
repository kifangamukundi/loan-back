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

type RoadController struct {
	RoadModel *models.RoadModel
}

func NewRoadController(roadModel *models.RoadModel) *RoadController {
	return &RoadController{RoadModel: roadModel}
}

// Create a road with body {roadName, country_id}
func (ctrl *RoadController) CreateRoadController(c *gin.Context) {
	var req bindings.CreateRoadRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	road := models.Road{
		RoadName: parameters.TrimWhitespace(req.RoadName),
		VillageID:  uint(req.VillageID),
	}

	if err := ctrl.RoadModel.CreateRoad(&road); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of roads with pagination
func (ctrl *RoadController) GetRoadsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roads, totalCount, count, err := ctrl.RoadModel.GetRoads(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching roads: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedRoads := transformations.Transform(roads, fieldNames,
		func(road models.Road) interface{} { return road.ID },
		func(road models.Road) interface{} { return road.RoadName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedRoads)
}

// Get a road by param {id}
func (ctrl *RoadController) GetRoadByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	road, err := ctrl.RoadModel.GetRoadByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Road not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, road)
}

// Update a road with body {roadName, country_id} and param {id}
func (ctrl *RoadController) UpdateRoadController(c *gin.Context) {
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

	var req bindings.UpdateRoadRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.RoadName)

	road, err := ctrl.RoadModel.UpdateRoad(idInt, name, uint(req.VillageID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating road: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, road)
}

// Delete road by param {id}
func (ctrl *RoadController) DeleteRoadController(c *gin.Context) {
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

	err = ctrl.RoadModel.DeleteRoad(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting road: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all roads formatted in {id and title}
func (ctrl *RoadController) GetAllRoadsController(c *gin.Context) {
	roads, err := ctrl.RoadModel.GetAllRoads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting roads: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedRoads := transformations.Transform(roads, fieldNames,
		func(road models.Road) interface{} { return road.ID },
		func(road models.Road) interface{} { return road.RoadName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedRoads)
}

func (ctrl *RoadController) GetRoadPlotsController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	road, err := ctrl.RoadModel.GetRoadPlotsByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Road not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedPlots := transformations.Transform(road.Plots, fieldNames,
		func(plot models.Plot) interface{} { return plot.ID },
		func(plot models.Plot) interface{} { return plot.PlotName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedPlots)
}

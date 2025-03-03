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

type VillageController struct {
	VillageModel *models.VillageModel
}

func NewVillageController(villageModel *models.VillageModel) *VillageController {
	return &VillageController{VillageModel: villageModel}
}

// Create a village with body {villageName, country_id}
func (ctrl *VillageController) CreateVillageController(c *gin.Context) {
	var req bindings.CreateVillageRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	village := models.Village{
		VillageName: parameters.TrimWhitespace(req.VillageName),
		SubLocationID:  uint(req.SubLocationID),
	}

	if err := ctrl.VillageModel.CreateVillage(&village); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of villages with pagination
func (ctrl *VillageController) GetVillagesController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	villages, totalCount, count, err := ctrl.VillageModel.GetVillages(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching villages: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedVillages := transformations.Transform(villages, fieldNames,
		func(village models.Village) interface{} { return village.ID },
		func(village models.Village) interface{} { return village.VillageName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedVillages)
}

// Get a village by param {id}
func (ctrl *VillageController) GetVillageByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	village, err := ctrl.VillageModel.GetVillageByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Village not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, village)
}

// Update a village with body {villageName, country_id} and param {id}
func (ctrl *VillageController) UpdateVillageController(c *gin.Context) {
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

	var req bindings.UpdateVillageRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.VillageName)

	village, err := ctrl.VillageModel.UpdateVillage(idInt, name, uint(req.SubLocationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating village: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, village)
}

// Delete village by param {id}
func (ctrl *VillageController) DeleteVillageController(c *gin.Context) {
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

	err = ctrl.VillageModel.DeleteVillage(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting village: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all villages formatted in {id and title}
func (ctrl *VillageController) GetAllVillagesController(c *gin.Context) {
	villages, err := ctrl.VillageModel.GetAllVillages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting villages: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedVillages := transformations.Transform(villages, fieldNames,
		func(village models.Village) interface{} { return village.ID },
		func(village models.Village) interface{} { return village.VillageName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedVillages)
}

func (ctrl *VillageController) GetVillageRoadsController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	village, err := ctrl.VillageModel.GetVillageRoadsByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Village not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedRoads := transformations.Transform(village.Roads, fieldNames,
		func(road models.Road) interface{} { return road.ID },
		func(road models.Road) interface{} { return road.RoadName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedRoads)
}

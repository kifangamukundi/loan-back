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

type SubCountyController struct {
	SubCountyModel *models.SubCountyModel
}

func NewSubCountyController(subCountyModel *models.SubCountyModel) *SubCountyController {
	return &SubCountyController{SubCountyModel: subCountyModel}
}

// Create a subCounty with body {subCountyName, country_id}
func (ctrl *SubCountyController) CreateSubCountyController(c *gin.Context) {
	var req bindings.CreateSubCountyRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	subCounty := models.SubCounty{
		SubCountyName: parameters.TrimWhitespace(req.SubCountyName),
		CountyID:     uint(req.CountyID),
	}

	if err := ctrl.SubCountyModel.CreateSubCounty(&subCounty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

// Get a list of subCountys with pagination
func (ctrl *SubCountyController) GetSubCountysController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subCountys, totalCount, count, err := ctrl.SubCountyModel.GetSubCountys(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching subCountys: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedSubCountys := transformations.Transform(subCountys, fieldNames,
		func(subCounty models.SubCounty) interface{} { return subCounty.ID },
		func(subCounty models.SubCounty) interface{} { return subCounty.SubCountyName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedSubCountys)
}

// Get a subCounty by param {id}
func (ctrl *SubCountyController) GetSubCountyByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	subCounty, err := ctrl.SubCountyModel.GetSubCountyByField("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubCounty not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, subCounty)
}

// Update a subCounty with body {subCountyName, country_id} and param {id}
func (ctrl *SubCountyController) UpdateSubCountyController(c *gin.Context) {
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

	var req bindings.UpdateSubCountyRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	name := parameters.TrimWhitespace(req.SubCountyName)

	subCounty, err := ctrl.SubCountyModel.UpdateSubCounty(idInt, name, uint(req.CountyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating subCounty: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, subCounty)
}

// Delete subCounty by param {id}
func (ctrl *SubCountyController) DeleteSubCountyController(c *gin.Context) {
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

	err = ctrl.SubCountyModel.DeleteSubCounty(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting subCounty: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

// Get all subCountys formatted in {id and title}
func (ctrl *SubCountyController) GetAllSubCountysController(c *gin.Context) {
	subCountys, err := ctrl.SubCountyModel.GetAllSubCountys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting subCountys: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedSubCountys := transformations.Transform(subCountys, fieldNames,
		func(subCounty models.SubCounty) interface{} { return subCounty.ID },
		func(subCounty models.SubCounty) interface{} { return subCounty.SubCountyName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedSubCountys)
}

func (ctrl *SubCountyController) GetSubCountyWardsController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	subCounty, err := ctrl.SubCountyModel.GetSubCountyWardsByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubCounty not found"})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedSubCounties := transformations.Transform(subCounty.Wards, fieldNames,
		func(ward models.Ward) interface{} { return ward.ID },
		func(ward models.Ward) interface{} { return ward.WardName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedSubCounties)
}

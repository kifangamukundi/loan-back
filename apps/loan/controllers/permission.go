package controllers

import (
	"fmt"
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

type PermissionController struct {
	PermissionModel *models.PermissionModel
}

func NewPermissionController(permissionModel *models.PermissionModel) *PermissionController {
	return &PermissionController{PermissionModel: permissionModel}
}

func (ctrl *PermissionController) CreatePermissionController(c *gin.Context) {
	var req bindings.CreatePermissionRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	if _, err := ctrl.PermissionModel.CreatePermission(req.PermissionName, req.Roles); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "some roles do not exist" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

func (ctrl *PermissionController) GetPermissionsController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissions, totalCount, count, err := ctrl.PermissionModel.GetPermissions(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching permissions: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedPermissions := transformations.Transform(permissions, fieldNames,
		func(permission models.Permission) interface{} { return permission.ID },
		func(permission models.Permission) interface{} { return permission.PermissionName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedPermissions)
}

func (ctrl *PermissionController) GetAllPermissionsController(c *gin.Context) {
	permissions, err := ctrl.PermissionModel.GetAllPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting permissions: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedPermissions := transformations.Transform(permissions, fieldNames,
		func(permission models.Permission) interface{} { return permission.ID },
		func(permission models.Permission) interface{} { return permission.PermissionName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedPermissions)
}

func (ctrl *PermissionController) GetPermissionByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	permission, err := ctrl.PermissionModel.GetPermissionByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, permission)
}

func (ctrl *PermissionController) UpdatePermissionController(c *gin.Context) {
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

	var req bindings.UpdatePermissionRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	permission, err := ctrl.PermissionModel.UpdatePermission(idInt, req.PermissionName, req.Roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating permission: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, permission)
}

func (ctrl *PermissionController) DeletePermissionController(c *gin.Context) {
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

	err = ctrl.PermissionModel.DeletePermission(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting permission: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

func (ctrl *PermissionController) GetAllPermissionsByPermissionsController(c *gin.Context) {
	permissions, err := ctrl.PermissionModel.GetAllPermissionsByRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting permissions: " + err.Error()})
		return
	}

	fieldNames := []string{"name", "value"}

	transformedPermissions := transformations.Transform(permissions, fieldNames,
		func(role models.Permission) interface{} { return role.PermissionName },
		func(role models.Permission) interface{} { return len(role.Roles) },
	)

	binders.ReturnJSONGeneralResponse(c, transformedPermissions)
}

func (ctrl *PermissionController) CountPermissionsController(c *gin.Context) {
	conditions := map[string]interface{}{}

	count, err := ctrl.PermissionModel.CountPermissions(conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to count permissions: %v", err)})
		return
	}

	transformedCount := map[string]interface{}{
		"count": count,
	}

	binders.ReturnJSONGeneralResponse(c, transformedCount)
}

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

type RoleController struct {
	RoleModel *models.RoleModel
}

func NewRoleController(roleModel *models.RoleModel) *RoleController {
	return &RoleController{RoleModel: roleModel}
}

func (ctrl *RoleController) CreateRoleController(c *gin.Context) {
	var req bindings.CreateRoleRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	if _, err := ctrl.RoleModel.CreateRole(req.RoleName, req.Permissions); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "some permissions do not exist" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	binders.ReturnJSONCreatedGenericResponse(c)
}

func (ctrl *RoleController) GetRolesController(c *gin.Context) {
	page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteria, err := queryparams.ExtractPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roles, totalCount, count, err := ctrl.RoleModel.GetRoles(skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching roles: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedRoles := transformations.Transform(roles, fieldNames,
		func(role models.Role) interface{} { return role.ID },
		func(role models.Role) interface{} { return role.RoleName },
	)

	binders.ReturnJSONPaginateResponse(c, page, limit, int(totalCount), int(count), transformedRoles)
}

func (ctrl *RoleController) GetAllRolesController(c *gin.Context) {
	roles, err := ctrl.RoleModel.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting roles: " + err.Error()})
		return
	}

	fieldNames := []string{
		"id",
		"title",
	}

	transformedRoles := transformations.Transform(roles, fieldNames,
		func(role models.Role) interface{} { return role.ID },
		func(role models.Role) interface{} { return role.RoleName },
	)

	binders.ReturnJSONGeneralResponse(c, transformedRoles)
}

func (ctrl *RoleController) GetRoleByIdController(c *gin.Context) {
	id, valid := parameters.ConvertParamToValidID(c, "id")
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	role, err := ctrl.RoleModel.GetRoleByFieldPreloaded("id", string(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	binders.ReturnJSONGeneralResponse(c, role)
}

func (ctrl *RoleController) UpdateRoleController(c *gin.Context) {
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

	var req bindings.UpdateRoleRequest
	if !binders.ValidateBindJSONRequest(c, &req) {
		return
	}

	role, err := ctrl.RoleModel.UpdateRole(idInt, req.RoleName, req.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating role: " + err.Error()})
		return
	}

	binders.ReturnJSONGeneralResponse(c, role)
}

func (ctrl *RoleController) DeleteRoleController(c *gin.Context) {
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

	err = ctrl.RoleModel.DeleteRole(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting role: " + err.Error()})
		return
	}

	binders.ReturnJSONOkayGenericResponse(c)
}

func (ctrl *RoleController) GetAllUsersByRoleController(c *gin.Context) {
	roles, err := ctrl.RoleModel.GetAllUsersByRole()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting roles: " + err.Error()})
		return
	}

	fieldNames := []string{"name", "users"}

	transformedRoles := transformations.Transform(roles, fieldNames,
		func(role models.Role) interface{} { return role.RoleName },
		func(role models.Role) interface{} { return len(role.Users) },
	)

	binders.ReturnJSONGeneralResponse(c, transformedRoles)
}

func (ctrl *RoleController) CountRolesController(c *gin.Context) {
	conditions := map[string]interface{}{}

	count, err := ctrl.RoleModel.CountRoles(conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to count roles: %v", err)})
		return
	}

	transformedCount := map[string]interface{}{
		"count": count,
	}

	binders.ReturnJSONGeneralResponse(c, transformedCount)
}

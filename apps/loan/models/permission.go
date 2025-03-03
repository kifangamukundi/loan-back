package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Permission struct {
	ID             uint   `gorm:"primaryKey"`
	PermissionName string `gorm:"unique;not null;index"`
	Roles          []Role `gorm:"many2many:role_permissions;"`
}

type PermissionModel struct {
	Service services.Service
}

func NewPermissionModel(service services.Service) *PermissionModel {
	return &PermissionModel{Service: service}
}

// Uses the services interface which in turn uses the repository interface
func (m *PermissionModel) CreatePermission(permissionName string, roleIds []int) (Permission, error) {
	var roles []Role
	if len(roleIds) > 0 {
		rolesMap := map[string]interface{}{
			"id": roleIds,
		}

		result, err := m.Service.GetEntitiesByFields(&[]Role{}, rolesMap)
		if err != nil {
			return Permission{}, fmt.Errorf("error fetching roles: %v", err)
		}

		rolesSlice, ok := result.(*[]Role)
		if !ok {
			return Permission{}, fmt.Errorf("unexpected type for roles result: %T", result)
		}

		roles = append(roles, *rolesSlice...)
	}

	name := parameters.TrimWhitespace(permissionName)

	permission := Permission{
		PermissionName: name,
		Roles:          roles,
	}

	if err := m.Service.CreateEntity(&permission); err != nil {
		return Permission{}, fmt.Errorf("failed to create permission: %v", err)
	}

	return permission, nil
}

func (m *PermissionModel) UpdatePermission(id int, PermissionName string, roleIds []int) (Permission, error) {
	var permission Permission
	result, err := m.Service.GetEntityByID(&permission, uint(id))
	if err != nil {
		return Permission{}, fmt.Errorf("permission not found: %v", err)
	}

	fetchedPermission, ok := result.(*Permission)
	if !ok {
		return Permission{}, fmt.Errorf("unexpected type for permission result: %T", result)
	}

	fetchedPermission.PermissionName = parameters.TrimWhitespace(PermissionName)

	var roles []Role
	if len(roleIds) > 0 {
		rolesMap := map[string]interface{}{
			"id": roleIds,
		}

		roleResult, err := m.Service.GetEntitiesByFields(&[]Role{}, rolesMap)
		if err != nil {
			return Permission{}, fmt.Errorf("error fetching roles: %v", err)
		}

		rolesSlice, ok := roleResult.(*[]Role)
		if !ok {
			return Permission{}, fmt.Errorf("unexpected type for roles result: %T", roleResult)
		}

		roles = append(roles, *rolesSlice...)
	}

	if err := m.Service.EntityClearAssociation(fetchedPermission, "Roles"); err != nil {
		return Permission{}, fmt.Errorf("failed to clear old roles: %v", err)
	}

	fetchedPermission.Roles = roles

	if err := m.Service.UpdateEntity(fetchedPermission); err != nil {
		return Permission{}, fmt.Errorf("failed to update permission: %v", err)
	}

	return *fetchedPermission, nil
}

func (m *PermissionModel) GetPermissions(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Permission, int64, int64, error) {
	searchColumns := []string{"permission_name"}

	preloads := []string{"Roles"}

	permissionsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Permission{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get permissions: %v", err)
	}

	var permissions []Permission
	for _, permission := range permissionsResult {
		if c, ok := permission.(*Permission); ok {
			permissions = append(permissions, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", permission)
		}
	}

	return permissions, totalCount, filteredCount, nil
}

func (m *PermissionModel) GetAllPermissions() ([]Permission, error) {
	var permissions []Permission

	result, err := m.Service.GetAllEntities(&permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %v", err)
	}

	permissionsPtr, ok := result.(*[]Permission)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *permissionsPtr, nil
}

func (m *PermissionModel) GetPermissionByField(field, value string) (*Permission, error) {
	var permission Permission

	result, err := m.Service.GetEntityByField(field, value, &permission)
	if err != nil {
		log.Printf("Error fetching permission by %s: %v", field, err)
		return nil, err
	}

	return result.(*Permission), nil
}

func (m *PermissionModel) GetPermissionByFieldPreloaded(field, value string) (*Permission, error) {
	var permission Permission

	preloads := []string{"Roles"}

	result, err := m.Service.GetEntityByFieldWithPreload(&permission, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching permission with roles by %s: %v", field, err)
		return nil, err
	}

	permissionPtr, ok := result.(*Permission)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return permissionPtr, nil
}

func (m *PermissionModel) DeletePermission(id uint) error {
	permission := &Permission{}

	_, err := m.Service.GetEntityByID(permission, id)
	if err != nil {
		return fmt.Errorf("permission not found: %v", err)
	}

	associations := []string{"Roles"}

	if err := m.Service.HardDeleteEntity(permission, id, "permission", associations...); err != nil {
		return fmt.Errorf("failed to delete permission: %v", err)
	}

	return nil
}

func (m *PermissionModel) GetAllPermissionsByRoles() ([]Permission, error) {
	preloads := []string{"Roles"}

	var permissions []Permission

	result, err := m.Service.GetAllEntitiesWithPreload(&permissions, preloads...)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %v", err)
	}

	permissionsPtr, ok := result.(*[]Permission)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *permissionsPtr, nil
}

func (m *PermissionModel) CountPermissions(conditions map[string]interface{}) (int64, error) {
	count, err := m.Service.CountEntities(&Permission{}, conditions)
	if err != nil {
		log.Printf("Error counting permissions: %v", err)
		return 0, fmt.Errorf("failed to count permissions: %v", err)
	}
	return count, nil
}
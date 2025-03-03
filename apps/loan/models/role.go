package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Role struct {
	ID          uint         `gorm:"primaryKey"`
	RoleName    string       `gorm:"unique;not null;index"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Users       []User       `gorm:"many2many:user_roles"`
}

type RoleModel struct {
	Service services.Service
}

func NewRoleModel(service services.Service) *RoleModel {
	return &RoleModel{Service: service}
}

// Uses the services interface which in turn uses the repository interface

func (m *RoleModel) CreateRole(RoleName string, permissionIDs []int) (Role, error) {
	var permissions []Permission
	if len(permissionIDs) > 0 {
		permissionsMap := map[string]interface{}{
			"id": permissionIDs,
		}

		result, err := m.Service.GetEntitiesByFields(&[]Permission{}, permissionsMap)
		if err != nil {
			return Role{}, fmt.Errorf("error fetching permissions: %v", err)
		}

		permissionsSlice, ok := result.(*[]Permission)
		if !ok {
			return Role{}, fmt.Errorf("unexpected type for permissions result: %T", result)
		}

		permissions = append(permissions, *permissionsSlice...)
	}

	name := parameters.TrimWhitespace(RoleName)

	role := Role{
		RoleName:    name,
		Permissions: permissions,
	}

	if err := m.Service.CreateEntity(&role); err != nil {
		return Role{}, fmt.Errorf("failed to create role: %v", err)
	}

	return role, nil
}

func (m *RoleModel) UpdateRole(id int, RoleName string, permissionIDs []int) (Role, error) {
	var role Role
	result, err := m.Service.GetEntityByID(&role, uint(id))
	if err != nil {
		return Role{}, fmt.Errorf("role not found: %v", err)
	}

	fetchedRole, ok := result.(*Role)
	if !ok {
		return Role{}, fmt.Errorf("unexpected type for role result: %T", result)
	}

	fetchedRole.RoleName = parameters.TrimWhitespace(RoleName)

	var permissions []Permission
	if len(permissionIDs) > 0 {
		permissionsMap := map[string]interface{}{
			"id": permissionIDs,
		}

		postResult, err := m.Service.GetEntitiesByFields(&[]Permission{}, permissionsMap)
		if err != nil {
			return Role{}, fmt.Errorf("error fetching permissions: %v", err)
		}

		permissionsSlice, ok := postResult.(*[]Permission)
		if !ok {
			return Role{}, fmt.Errorf("unexpected type for permissions result: %T", postResult)
		}

		permissions = append(permissions, *permissionsSlice...)
	}

	if err := m.Service.EntityClearAssociation(fetchedRole, "Permissions"); err != nil {
		return Role{}, fmt.Errorf("failed to clear old permissions: %v", err)
	}

	fetchedRole.Permissions = permissions

	if err := m.Service.UpdateEntity(fetchedRole); err != nil {
		return Role{}, fmt.Errorf("failed to update role: %v", err)
	}

	return *fetchedRole, nil
}

func (m *RoleModel) GetRoles(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Role, int64, int64, error) {
	searchColumns := []string{"role_name"}

	preloads := []string{"Permissions"}

	rolesResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Role{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get roles: %v", err)
	}

	var roles []Role
	for _, role := range rolesResult {
		if c, ok := role.(*Role); ok {
			roles = append(roles, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", role)
		}
	}

	return roles, totalCount, filteredCount, nil
}

func (m *RoleModel) GetAllRoles() ([]Role, error) {
	var roles []Role

	result, err := m.Service.GetAllEntities(&roles)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %v", err)
	}

	rolesPtr, ok := result.(*[]Role)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *rolesPtr, nil
}

func (m *RoleModel) GetRoleByField(field, value string) (*Role, error) {
	var role Role

	result, err := m.Service.GetEntityByField(field, value, &role)
	if err != nil {
		log.Printf("Error fetching role by %s: %v", field, err)
		return nil, err
	}

	return result.(*Role), nil
}

func (m *RoleModel) GetRoleByFieldPreloaded(field, value string) (*Role, error) {
	var role Role

	preloads := []string{"Permissions"}

	result, err := m.Service.GetEntityByFieldWithPreload(&role, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching role with permissions by %s: %v", field, err)
		return nil, err
	}

	rolePtr, ok := result.(*Role)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return rolePtr, nil
}

func (m *RoleModel) DeleteRole(id uint) error {
	role := &Role{}

	_, err := m.Service.GetEntityByID(role, id)
	if err != nil {
		return fmt.Errorf("role not found: %v", err)
	}

	associations := []string{"Permissions", "Users"}

	if err := m.Service.HardDeleteEntity(role, id, "role", associations...); err != nil {
		return fmt.Errorf("failed to delete role: %v", err)
	}

	return nil
}

func (m *RoleModel) GetAllUsersByRole() ([]Role, error) {
	preloads := []string{"Users"}

	var roles []Role

	result, err := m.Service.GetAllEntitiesWithPreload(&roles, preloads...)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %v", err)
	}

	rolesPtr, ok := result.(*[]Role)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *rolesPtr, nil
}

func (m *RoleModel) CountRoles(conditions map[string]interface{}) (int64, error) {
	count, err := m.Service.CountEntities(&Role{}, conditions)
	if err != nil {
		log.Printf("Error counting roles: %v", err)
		return 0, fmt.Errorf("failed to count roles: %v", err)
	}
	return count, nil
}

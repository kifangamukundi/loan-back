package models

import (
	"fmt"
	"log"
	"time"

	"github.com/kifangamukundi/gm/loan/services"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	FirstName    string `gorm:"not null;index"`
	LastName     string `gorm:"not null;index"`
	Email        string `gorm:"unique;not null;index"`
	MobileNumber string `gorm:"unique;not null;index"`
	Password     string `gorm:"not null"`
	IsActive     bool   `gorm:"default:false"`
	IsLocked     bool   `gorm:"default:false"`

	AccountActivationToken  string    `gorm:"index"`
	AccountActivationExpire time.Time `gorm:"not null"`
	ResetPasswordToken      string    `gorm:"index"`
	ResetPasswordExpire     time.Time `gorm:"not null"`
	ResetRequestCount       int       `gorm:"default:0"`
	LastResetRequestAt      time.Time

	Roles []Role `gorm:"many2many:user_roles"`
	Agent *Agent `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time `gorm:"index"`
}

type UserModel struct {
	Service services.Service
}

func NewUserModel(service services.Service) *UserModel {
	return &UserModel{Service: service}
}

func (m *UserModel) GetMonthlyUserCounts() ([]map[string]interface{}, error) {
    return m.Service.GetMonthlyEntityCounts(&User{})
}

// Uses the services interface which in turn uses the repository interface

func (m *UserModel) CreateUser(user *User) error {
	if err := m.Service.CreateEntity(user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (m *UserModel) CountUsers(conditions map[string]interface{}) (int64, error) {
	count, err := m.Service.CountEntities(&User{}, conditions)
	if err != nil {
		log.Printf("Error counting users: %v", err)
		return 0, fmt.Errorf("failed to count users: %v", err)
	}
	return count, nil
}

func (m *UserModel) CreateUserAgent(user *User) error {
	roleNames := []string{"Agent"}

	if len(roleNames) > 0 {
		rolesMap := map[string]interface{}{
			"role_name": roleNames,
		}

		result, err := m.Service.GetEntitiesByFields(&[]Role{}, rolesMap)
		if err != nil {
			return fmt.Errorf("error fetching roles: %v", err)
		}

		rolesSlice, ok := result.(*[]Role)
		if !ok {
			return fmt.Errorf("unexpected type for roles result: %T", result)
		}

		user.Roles = *rolesSlice
	}

	if err := m.Service.CreateEntity(user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (m *UserModel) CreateUserOfficer(user *User) error {
	roleNames := []string{"Officer"}

	if len(roleNames) > 0 {
		rolesMap := map[string]interface{}{
			"role_name": roleNames,
		}

		result, err := m.Service.GetEntitiesByFields(&[]Role{}, rolesMap)
		if err != nil {
			return fmt.Errorf("error fetching roles: %v", err)
		}

		rolesSlice, ok := result.(*[]Role)
		if !ok {
			return fmt.Errorf("unexpected type for roles result: %T", result)
		}

		user.Roles = *rolesSlice
	}

	if err := m.Service.CreateEntity(user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (m *UserModel) CreateUserMember(user *User) error {
	roleNames := []string{"Member"}

	if len(roleNames) > 0 {
		rolesMap := map[string]interface{}{
			"role_name": roleNames,
		}

		result, err := m.Service.GetEntitiesByFields(&[]Role{}, rolesMap)
		if err != nil {
			return fmt.Errorf("error fetching roles: %v", err)
		}

		rolesSlice, ok := result.(*[]Role)
		if !ok {
			return fmt.Errorf("unexpected type for roles result: %T", result)
		}

		user.Roles = *rolesSlice
	}

	if err := m.Service.CreateEntity(user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (m *UserModel) GetUserByField(field, value string) (*User, error) {
	var user User

	result, err := m.Service.GetEntityByField(field, value, &user)
	if err != nil {
		log.Printf("Error fetching user by %s: %v", field, err)
		return nil, err
	}

	return result.(*User), nil
}

func (m *UserModel) UpdateUserGeneric(user *User) error {
	if err := m.Service.UpdateEntity(user); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func (m *UserModel) DeleteUserPermanent(id uint) error {
	user := &User{}

	_, err := m.Service.GetEntityByID(user, id)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Todo: Also delete user posts, comments and likes

	associations := []string{"Roles", "Posts", "Comments", "Likes"}

	if err := m.Service.HardDeleteEntity(user, id, "user", associations...); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

func (m *UserModel) DeleteUserTemporary(id uint) error {
	user := &User{}

	_, err := m.Service.GetEntityByID(user, id)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	associations := []string{}

	if err := m.Service.SoftDeleteEntity(user, id, "user", associations...); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

func (m *UserModel) GetUsers(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]User, int64, int64, error) {
	searchColumns := []string{"first_name", "last_name", "email", "mobile_number"}

	preloads := []string{"Roles"}

	usersResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&User{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get users: %v", err)
	}

	var users []User
	for _, user := range usersResult {
		if c, ok := user.(*User); ok {
			users = append(users, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", user)
		}
	}

	return users, totalCount, filteredCount, nil
}

func (m *UserModel) GetUserByFieldPreloaded(field, value string) (*User, error) {
	var user User

	preloads := []string{"Roles", "Agent"}

	result, err := m.Service.GetEntityByFieldWithPreload(&user, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching user with roles by %s: %v", field, err)
		return nil, err
	}

	userPtr, ok := result.(*User)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return userPtr, nil
}

func (m *UserModel) UpdateUser(id int, isActive, isLocked bool, roleIDs []int) (User, error) {
	var user User
	result, err := m.Service.GetEntityByID(&user, uint(id))
	if err != nil {
		return User{}, fmt.Errorf("user not found: %v", err)
	}

	fetchedUser, ok := result.(*User)
	if !ok {
		return User{}, fmt.Errorf("unexpected type for user result: %T", result)
	}

	fetchedUser.IsActive = isActive
	fetchedUser.IsLocked = isLocked

	var roles []Role
	if len(roleIDs) > 0 {
		rolesMap := map[string]interface{}{
			"id": roleIDs,
		}

		postResult, err := m.Service.GetEntitiesByFields(&[]Role{}, rolesMap)
		if err != nil {
			return User{}, fmt.Errorf("error fetching roles: %v", err)
		}

		rolesSlice, ok := postResult.(*[]Role)
		if !ok {
			return User{}, fmt.Errorf("unexpected type for roles result: %T", postResult)
		}

		roles = append(roles, *rolesSlice...)
	}

	if err := m.Service.EntityClearAssociation(fetchedUser, "Roles"); err != nil {
		return User{}, fmt.Errorf("failed to clear old roles: %v", err)
	}

	fetchedUser.Roles = roles

	if err := m.Service.UpdateEntity(fetchedUser); err != nil {
		return User{}, fmt.Errorf("failed to update user: %v", err)
	}

	return *fetchedUser, nil
}

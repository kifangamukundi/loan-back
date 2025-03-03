package models

import (
	"fmt"
	"log"
	"time"

	"github.com/kifangamukundi/gm/loan/services"
)

type Group struct {
	ID        uint   `gorm:"primaryKey"`
	GroupName string `gorm:"unique;not null;index"`
	IsActive  bool   `gorm:"default:true"`

	AgentID uint  `gorm:"index"`
	Agent   Agent `gorm:"foreignKey:AgentID"`

	CountryID uint    `gorm:"index"`
	Country   Country `gorm:"foreignKey:CountryID"`

	RegionID uint   `gorm:"index"`
	Region   Region `gorm:"foreignKey:RegionID"`

	CityID uint `gorm:"index"`
	City   City `gorm:"foreignKey:CityID"`

	Members []Member `gorm:"many2many:group_members"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type GroupModel struct {
	Service services.Service
}

func NewGroupModel(service services.Service) *GroupModel {
	return &GroupModel{Service: service}
}

func (m *GroupModel) CreateGroup(group *Group) error {
	if err := m.Service.CreateEntity(group); err != nil {
		return fmt.Errorf("failed to create group: %v", err)
	}

	return nil
}

func (m *GroupModel) GetGroups(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Group, int64, int64, error) {
	searchColumns := []string{"group_name"}

	preloads := []string{}

	groupsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Group{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get groups: %v", err)
	}

	var groups []Group
	for _, agent := range groupsResult {
		if c, ok := agent.(*Group); ok {
			groups = append(groups, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", agent)
		}
	}

	return groups, totalCount, filteredCount, nil
}

func (m *GroupModel) GetGroupByField(field, value string) (*Group, error) {
	var group Group

	result, err := m.Service.GetEntityByField(field, value, &group)
	if err != nil {
		log.Printf("Error fetching group by %s: %v", field, err)
		return nil, err
	}

	return result.(*Group), nil
}

func (m *GroupModel) UpdateGroup(id int, isActive bool, agentID uint) (Group, error) {
	var group Group
	result, err := m.Service.GetEntityByID(&group, uint(id))
	if err != nil {
		return Group{}, fmt.Errorf("group not found: %v", err)
	}

	fetchedGroup, ok := result.(*Group)
	if !ok {
		return Group{}, fmt.Errorf("unexpected type for group result: %T", result)
	}

	fetchedGroup.IsActive = isActive
	fetchedGroup.AgentID = agentID

	if err := m.Service.UpdateEntity(fetchedGroup); err != nil {
		return Group{}, fmt.Errorf("failed to update group: %v", err)
	}

	return *fetchedGroup, nil
}

func (m *GroupModel) GetAgentGroups(agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Group, int64, int64, error) {
	searchColumns := []string{"user.first_name", "user.last_name"}

	preloads := []string{
		"Agent.User",
		"Members.User",
	}

	groupsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFilteredTest(&Group{}, agentId, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get groups: %v", err)
	}

	var groups []Group
	for _, group := range groupsResult {
		if c, ok := group.(*Group); ok {
			groups = append(groups, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", group)
		}
	}

	return groups, totalCount, filteredCount, nil
}

func (m *GroupModel) CountGroups(conditions map[string]interface{}) (int64, error) {
	count, err := m.Service.CountEntities(&Group{}, conditions)
	if err != nil {
		log.Printf("Error counting groups: %v", err)
		return 0, fmt.Errorf("failed to count groups: %v", err)
	}
	return count, nil
}

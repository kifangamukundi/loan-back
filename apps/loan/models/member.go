package models

import (
	"fmt"
	"log"
	"time"

	"github.com/kifangamukundi/gm/loan/services"
)

type Member struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"unique; index"`
	User   User `gorm:"foreignKey:UserID"`

	AgentID uint  `gorm:"index"`
	Agent   Agent `gorm:"foreignKey:AgentID"`

	IsActive  bool       `gorm:"default:false"`
	LastLogin *time.Time `gorm:"default:null"`

	CountryID uint    `gorm:"index"`
	Country   Country `gorm:"foreignKey:CountryID"`

	RegionID uint   `gorm:"index"`
	Region   Region `gorm:"foreignKey:RegionID"`

	CityID uint `gorm:"index"`
	City   City `gorm:"foreignKey:CityID"`

	Groups []Group `gorm:"many2many:group_members"`

	Loans []Loan `gorm:"foreignKey:MemberID;constraint:onDelete:CASCADE"`
}

type MemberModel struct {
	Service services.Service
}

func NewMemberModel(service services.Service) *MemberModel {
	return &MemberModel{Service: service}
}

func (m *MemberModel) CreateMember(agentId, userID uint, status bool, countryID, regionID, cityID uint, groupIDs []int) (Member, error) {
	var groups []Group
	if len(groupIDs) > 0 {
		groupsMap := map[string]interface{}{
			"id": groupIDs,
		}

		result, err := m.Service.GetEntitiesByFields(&[]Group{}, groupsMap)
		if err != nil {
			return Member{}, fmt.Errorf("error fetching groups: %v", err)
		}

		groupsSlice, ok := result.(*[]Group)
		if !ok {
			return Member{}, fmt.Errorf("unexpected type for groups result: %T", result)
		}

		groups = append(groups, *groupsSlice...)
	}

	member := Member{
		UserID:    userID,
		AgentID:   agentId,
		IsActive:  status,
		CountryID: countryID,
		RegionID:  regionID,
		CityID:    cityID,
		Groups:    groups,
	}

	if err := m.Service.CreateEntity(&member); err != nil {
		return Member{}, fmt.Errorf("failed to create member: %v", err)
	}

	return member, nil
}

func (m *MemberModel) GetGroupMembers(groupId, agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Member, int64, int64, error) {
	searchColumns := []string{"user.first_name", "user.last_name"}

	preloads := []string{
		"User",
		"Agent.User",
		"Groups.Members",
	}

	membersResult, totalCount, filteredCount, err := m.Service.GetEntitiesFilteredTest2(&Member{}, groupId, agentId, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get members: %v", err)
	}

	var members []Member
	for _, member := range membersResult {
		if c, ok := member.(*Member); ok {
			members = append(members, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", member)
		}
	}

	return members, totalCount, filteredCount, nil
}

func (m *MemberModel) GetMemberByField(field, value string) (*Member, error) {
	var member Member

	result, err := m.Service.GetEntityByField(field, value, &member)
	if err != nil {
		log.Printf("Error fetching member by %s: %v", field, err)
		return nil, err
	}

	return result.(*Member), nil
}

func (m *MemberModel) GetMemberByFieldPreloaded(field, value string) (*Member, error) {
	var member Member

	preloads := []string{"User"}

	result, err := m.Service.GetEntityByFieldWithPreload(&member, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching member with user info by %s: %v", field, err)
		return nil, err
	}

	memberPtr, ok := result.(*Member)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return memberPtr, nil
}

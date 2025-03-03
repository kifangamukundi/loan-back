package models

import (
	"fmt"
	"log"
	"time"

	"github.com/kifangamukundi/gm/loan/services"
)

type Officer struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"unique; index"`
	User   User `gorm:"foreignKey:UserID"`

	IsActive  bool       `gorm:"default:false"`
	LastLogin *time.Time `gorm:"default:null"`

	CountryID uint    `gorm:"index"`
	Country   Country `gorm:"foreignKey:CountryID"`

	RegionID uint   `gorm:"index"`
	Region   Region `gorm:"foreignKey:RegionID"`

	CityID uint `gorm:"index"`
	City   City `gorm:"foreignKey:CityID"`

	Loans []Loan `gorm:"foreignKey:OfficerID"`
}

type OfficerModel struct {
	Service services.Service
}

func NewOfficerModel(service services.Service) *OfficerModel {
	return &OfficerModel{Service: service}
}

func (m *OfficerModel) CreateOfficer(officer *Officer) error {
	if err := m.Service.CreateEntity(officer); err != nil {
		return fmt.Errorf("failed to create officer: %v", err)
	}

	return nil
}

func (m *OfficerModel) CountOfficers(conditions map[string]interface{}) (int64, error) {
	count, err := m.Service.CountEntities(&Officer{}, conditions)
	if err != nil {
		log.Printf("Error counting officers: %v", err)
		return 0, fmt.Errorf("failed to count officers: %v", err)
	}
	return count, nil
}

func (m *OfficerModel) GetOfficers(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Officer, int64, int64, error) {
	searchColumns := []string{"user.first_name", "user.last_name"}

	preloads := []string{"User"}

	officersResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Officer{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get officers: %v", err)
	}

	var officers []Officer
	for _, agent := range officersResult {
		if c, ok := agent.(*Officer); ok {
			officers = append(officers, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", agent)
		}
	}

	return officers, totalCount, filteredCount, nil
}

func (m *OfficerModel) GetAllOfficers() ([]Officer, error) {
	preloads := []string{"User"}

	var officers []Officer

	result, err := m.Service.GetAllEntitiesWithPreload(&officers, preloads...)
	if err != nil {
		return nil, fmt.Errorf("failed to get officers: %v", err)
	}

	officersPtr, ok := result.(*[]Officer)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *officersPtr, nil
}

func (m *OfficerModel) GetOfficerByFieldPreloaded(field, value string) (*Officer, error) {
	var officer Officer

	preloads := []string{"User"}

	result, err := m.Service.GetEntityByFieldWithPreload(&officer, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching officer with user info by %s: %v", field, err)
		return nil, err
	}

	officerPtr, ok := result.(*Officer)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return officerPtr, nil
}

func (m *OfficerModel) UpdateOfficer(id int, isActive bool) (Officer, error) {
	var officer Officer
	result, err := m.Service.GetEntityByID(&officer, uint(id))
	if err != nil {
		return Officer{}, fmt.Errorf("officer not found: %v", err)
	}

	fetchedOfficer, ok := result.(*Officer)
	if !ok {
		return Officer{}, fmt.Errorf("unexpected type for officer result: %T", result)
	}

	fetchedOfficer.IsActive = isActive

	if err := m.Service.UpdateEntity(fetchedOfficer); err != nil {
		return Officer{}, fmt.Errorf("failed to update officer: %v", err)
	}

	return *fetchedOfficer, nil
}

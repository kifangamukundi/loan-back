package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type SubCounty struct {
	ID uint `gorm:"primaryKey"`

	SubCountyName string `gorm:"unique;not null;index"`
	CountyID      uint   `gorm:"index"`

	County County `gorm:"foreignKey:CountyID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	Wards  []Ward `gorm:"foreignKey:SubCountyID;constraint:OnDelete:CASCADE"`
}

type SubCountyModel struct {
	Service services.Service
}

func NewSubCountyModel(service services.Service) *SubCountyModel {
	return &SubCountyModel{Service: service}
}

func (m *SubCountyModel) CreateSubCounty(subCounty *SubCounty) error {
	if err := m.Service.CreateEntity(subCounty); err != nil {
		return fmt.Errorf("failed to create subCounty: %v", err)
	}

	return nil
}

func (m *SubCountyModel) GetSubCountys(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]SubCounty, int64, int64, error) {
	searchColumns := []string{"sub_county_name"}

	preloads := []string{}

	subCountysResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&SubCounty{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get subCountys: %v", err)
	}

	var subCountys []SubCounty
	for _, subCounty := range subCountysResult {
		if c, ok := subCounty.(*SubCounty); ok {
			subCountys = append(subCountys, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", subCounty)
		}
	}

	return subCountys, totalCount, filteredCount, nil
}

func (m *SubCountyModel) GetSubCountyByField(field, value string) (*SubCounty, error) {
	var subCounty SubCounty

	result, err := m.Service.GetEntityByField(field, value, &subCounty)
	if err != nil {
		log.Printf("Error fetching subCounty by %s: %v", field, err)
		return nil, err
	}

	return result.(*SubCounty), nil
}

func (m *SubCountyModel) UpdateSubCounty(id int, subCountyName string, countyId uint) (SubCounty, error) {
	subCounty := &SubCounty{ID: uint(id)}

	_, err := m.Service.GetEntityByID(subCounty, uint(id))
	if err != nil {
		return SubCounty{}, fmt.Errorf("subCounty not found: %v", err)
	}

	subCounty.SubCountyName = parameters.TrimWhitespace(subCountyName)
	subCounty.CountyID = countyId

	if err := m.Service.UpdateEntity(subCounty); err != nil {
		return SubCounty{}, fmt.Errorf("failed to update subCounty: %v", err)
	}

	return *subCounty, nil
}

func (m *SubCountyModel) DeleteSubCounty(id uint) error {
	subCounty := &SubCounty{}

	_, err := m.Service.GetEntityByID(subCounty, id)
	if err != nil {
		return fmt.Errorf("subCounty not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(subCounty, id, "subCounty"); err != nil {
		return fmt.Errorf("failed to delete subCounty: %v", err)
	}

	return nil
}

func (m *SubCountyModel) GetAllSubCountys() ([]SubCounty, error) {
	var subCountys []SubCounty

	result, err := m.Service.GetAllEntities(&subCountys)
	if err != nil {
		return nil, fmt.Errorf("failed to get subCountys: %v", err)
	}

	subCountysPtr, ok := result.(*[]SubCounty)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *subCountysPtr, nil
}

func (m *SubCountyModel) GetSubCountyWardsByFieldPreloaded(field, value string) (*SubCounty, error) {
	var subCounty SubCounty

	preloads := []string{"Wards"}

	result, err := m.Service.GetEntityByFieldWithPreload(&subCounty, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching subCounty with counties by %s: %v", field, err)
		return nil, err
	}

	subCountyPtr, ok := result.(*SubCounty)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return subCountyPtr, nil
}

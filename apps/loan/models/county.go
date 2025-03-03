package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type County struct {
	ID uint `gorm:"primaryKey"`

	CountyName string `gorm:"unique;not null;index"`
	RegionID   uint   `gorm:"index"`

	Region      Region      `gorm:"foreignKey:RegionID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	SubCounties []SubCounty `gorm:"foreignKey:CountyID;constraint:OnDelete:CASCADE"`
}

type CountyModel struct {
	Service services.Service
}

func NewCountyModel(service services.Service) *CountyModel {
	return &CountyModel{Service: service}
}

func (m *CountyModel) CreateCounty(county *County) error {
	if err := m.Service.CreateEntity(county); err != nil {
		return fmt.Errorf("failed to create county: %v", err)
	}

	return nil
}

func (m *CountyModel) GetCountys(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]County, int64, int64, error) {
	searchColumns := []string{"county_name"}

	preloads := []string{}

	countysResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&County{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get countys: %v", err)
	}

	var countys []County
	for _, county := range countysResult {
		if c, ok := county.(*County); ok {
			countys = append(countys, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", county)
		}
	}

	return countys, totalCount, filteredCount, nil
}

func (m *CountyModel) GetCountyByField(field, value string) (*County, error) {
	var county County

	result, err := m.Service.GetEntityByField(field, value, &county)
	if err != nil {
		log.Printf("Error fetching county by %s: %v", field, err)
		return nil, err
	}

	return result.(*County), nil
}

func (m *CountyModel) UpdateCounty(id int, countyName string, regionId uint) (County, error) {
	county := &County{ID: uint(id)}

	_, err := m.Service.GetEntityByID(county, uint(id))
	if err != nil {
		return County{}, fmt.Errorf("county not found: %v", err)
	}

	county.CountyName = parameters.TrimWhitespace(countyName)
	county.RegionID = regionId

	if err := m.Service.UpdateEntity(county); err != nil {
		return County{}, fmt.Errorf("failed to update county: %v", err)
	}

	return *county, nil
}

func (m *CountyModel) DeleteCounty(id uint) error {
	county := &County{}

	_, err := m.Service.GetEntityByID(county, id)
	if err != nil {
		return fmt.Errorf("county not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(county, id, "county"); err != nil {
		return fmt.Errorf("failed to delete county: %v", err)
	}

	return nil
}

func (m *CountyModel) GetAllCountys() ([]County, error) {
	var countys []County

	result, err := m.Service.GetAllEntities(&countys)
	if err != nil {
		return nil, fmt.Errorf("failed to get countys: %v", err)
	}

	countysPtr, ok := result.(*[]County)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *countysPtr, nil
}

func (m *CountyModel) GetCountySubCountiesByFieldPreloaded(field, value string) (*County, error) {
	var county County

	preloads := []string{"SubCounties"}

	result, err := m.Service.GetEntityByFieldWithPreload(&county, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching county with counties by %s: %v", field, err)
		return nil, err
	}

	countyPtr, ok := result.(*County)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return countyPtr, nil
}

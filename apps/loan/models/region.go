package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Region struct {
	ID uint `gorm:"primaryKey"`

	RegionName string `gorm:"unique;not null;index"`
	CountryID  uint   `gorm:"index"`

	Country Country `gorm:"foreignKey:CountryID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	Counties  []County  `gorm:"foreignKey:RegionID;constraint:OnDelete:CASCADE"`
}

type RegionModel struct {
	Service services.Service
}

func NewRegionModel(service services.Service) *RegionModel {
	return &RegionModel{Service: service}
}

// Uses the services interface which in turn uses the repository interface
func (m *RegionModel) CreateRegion(region *Region) error {
	if err := m.Service.CreateEntity(region); err != nil {
		return fmt.Errorf("failed to create region: %v", err)
	}

	return nil
}

func (m *RegionModel) GetRegions(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Region, int64, int64, error) {
	searchColumns := []string{"region_name"}

	preloads := []string{}

	regionsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Region{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get regions: %v", err)
	}

	var regions []Region
	for _, region := range regionsResult {
		if c, ok := region.(*Region); ok {
			regions = append(regions, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", region)
		}
	}

	return regions, totalCount, filteredCount, nil
}

func (m *RegionModel) GetRegionByField(field, value string) (*Region, error) {
	var region Region

	result, err := m.Service.GetEntityByField(field, value, &region)
	if err != nil {
		log.Printf("Error fetching region by %s: %v", field, err)
		return nil, err
	}

	return result.(*Region), nil
}

func (m *RegionModel) UpdateRegion(id int, regionName string, countryId uint) (Region, error) {
	region := &Region{ID: uint(id)}

	_, err := m.Service.GetEntityByID(region, uint(id))
	if err != nil {
		return Region{}, fmt.Errorf("region not found: %v", err)
	}

	region.RegionName = parameters.TrimWhitespace(regionName)
	region.CountryID = countryId

	if err := m.Service.UpdateEntity(region); err != nil {
		return Region{}, fmt.Errorf("failed to update region: %v", err)
	}

	return *region, nil
}

func (m *RegionModel) DeleteRegion(id uint) error {
	region := &Region{}

	_, err := m.Service.GetEntityByID(region, id)
	if err != nil {
		return fmt.Errorf("region not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(region, id, "region"); err != nil {
		return fmt.Errorf("failed to delete region: %v", err)
	}

	return nil
}

func (m *RegionModel) GetAllRegions() ([]Region, error) {
	var regions []Region

	result, err := m.Service.GetAllEntities(&regions)
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %v", err)
	}

	regionsPtr, ok := result.(*[]Region)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *regionsPtr, nil
}

func (m *RegionModel) GetRegionCountiesByFieldPreloaded(field, value string) (*Region, error) {
	var region Region

	preloads := []string{"Counties"}

	result, err := m.Service.GetEntityByFieldWithPreload(&region, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching region with counties by %s: %v", field, err)
		return nil, err
	}

	regionPtr, ok := result.(*Region)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return regionPtr, nil
}

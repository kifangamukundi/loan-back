package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type SubLocation struct {
	ID uint `gorm:"primaryKey"`

	SubLocationName string `gorm:"unique;not null;index"`
	LocationID      uint   `gorm:"index"`

	Location Location  `gorm:"foreignKey:LocationID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	Villages []Village `gorm:"foreignKey:SubLocationID;constraint:OnDelete:CASCADE"`
}

type SubLocationModel struct {
	Service services.Service
}

func NewSubLocationModel(service services.Service) *SubLocationModel {
	return &SubLocationModel{Service: service}
}

func (m *SubLocationModel) CreateSubLocation(subLocation *SubLocation) error {
	if err := m.Service.CreateEntity(subLocation); err != nil {
		return fmt.Errorf("failed to create subLocation: %v", err)
	}

	return nil
}

func (m *SubLocationModel) GetSubLocations(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]SubLocation, int64, int64, error) {
	searchColumns := []string{"sub_location_name"}

	preloads := []string{}

	subLocationsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&SubLocation{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get subLocations: %v", err)
	}

	var subLocations []SubLocation
	for _, subLocation := range subLocationsResult {
		if c, ok := subLocation.(*SubLocation); ok {
			subLocations = append(subLocations, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", subLocation)
		}
	}

	return subLocations, totalCount, filteredCount, nil
}

func (m *SubLocationModel) GetSubLocationByField(field, value string) (*SubLocation, error) {
	var subLocation SubLocation

	result, err := m.Service.GetEntityByField(field, value, &subLocation)
	if err != nil {
		log.Printf("Error fetching subLocation by %s: %v", field, err)
		return nil, err
	}

	return result.(*SubLocation), nil
}

func (m *SubLocationModel) UpdateSubLocation(id int, subLocationName string, locationId uint) (SubLocation, error) {
	subLocation := &SubLocation{ID: uint(id)}

	_, err := m.Service.GetEntityByID(subLocation, uint(id))
	if err != nil {
		return SubLocation{}, fmt.Errorf("subLocation not found: %v", err)
	}

	subLocation.SubLocationName = parameters.TrimWhitespace(subLocationName)
	subLocation.LocationID = locationId

	if err := m.Service.UpdateEntity(subLocation); err != nil {
		return SubLocation{}, fmt.Errorf("failed to update subLocation: %v", err)
	}

	return *subLocation, nil
}

func (m *SubLocationModel) DeleteSubLocation(id uint) error {
	subLocation := &SubLocation{}

	_, err := m.Service.GetEntityByID(subLocation, id)
	if err != nil {
		return fmt.Errorf("subLocation not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(subLocation, id, "subLocation"); err != nil {
		return fmt.Errorf("failed to delete subLocation: %v", err)
	}

	return nil
}

func (m *SubLocationModel) GetAllSubLocations() ([]SubLocation, error) {
	var subLocations []SubLocation

	result, err := m.Service.GetAllEntities(&subLocations)
	if err != nil {
		return nil, fmt.Errorf("failed to get subLocations: %v", err)
	}

	subLocationsPtr, ok := result.(*[]SubLocation)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *subLocationsPtr, nil
}

func (m *SubLocationModel) GetSubLocationVillagesByFieldPreloaded(field, value string) (*SubLocation, error) {
	var subLocation SubLocation

	preloads := []string{"Villages"}

	result, err := m.Service.GetEntityByFieldWithPreload(&subLocation, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching subLocation with villages by %s: %v", field, err)
		return nil, err
	}

	subLocationPtr, ok := result.(*SubLocation)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return subLocationPtr, nil
}

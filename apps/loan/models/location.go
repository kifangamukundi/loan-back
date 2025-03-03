package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Location struct {
	ID uint `gorm:"primaryKey"`

	LocationName string `gorm:"unique;not null;index"`
	WardID       uint   `gorm:"index"`

	Ward         Ward          `gorm:"foreignKey:WardID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	SubLocations []SubLocation `gorm:"foreignKey:LocationID;constraint:OnDelete:CASCADE"`
}

type LocationModel struct {
	Service services.Service
}

func NewLocationModel(service services.Service) *LocationModel {
	return &LocationModel{Service: service}
}

func (m *LocationModel) CreateLocation(location *Location) error {
	if err := m.Service.CreateEntity(location); err != nil {
		return fmt.Errorf("failed to create location: %v", err)
	}

	return nil
}

func (m *LocationModel) GetLocations(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Location, int64, int64, error) {
	searchColumns := []string{"location_name"}

	preloads := []string{}

	locationsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Location{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get locations: %v", err)
	}

	var locations []Location
	for _, location := range locationsResult {
		if c, ok := location.(*Location); ok {
			locations = append(locations, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", location)
		}
	}

	return locations, totalCount, filteredCount, nil
}

func (m *LocationModel) GetLocationByField(field, value string) (*Location, error) {
	var location Location

	result, err := m.Service.GetEntityByField(field, value, &location)
	if err != nil {
		log.Printf("Error fetching location by %s: %v", field, err)
		return nil, err
	}

	return result.(*Location), nil
}

func (m *LocationModel) UpdateLocation(id int, locationName string, wardId uint) (Location, error) {
	location := &Location{ID: uint(id)}

	_, err := m.Service.GetEntityByID(location, uint(id))
	if err != nil {
		return Location{}, fmt.Errorf("location not found: %v", err)
	}

	location.LocationName = parameters.TrimWhitespace(locationName)
	location.WardID = wardId

	if err := m.Service.UpdateEntity(location); err != nil {
		return Location{}, fmt.Errorf("failed to update location: %v", err)
	}

	return *location, nil
}

func (m *LocationModel) DeleteLocation(id uint) error {
	location := &Location{}

	_, err := m.Service.GetEntityByID(location, id)
	if err != nil {
		return fmt.Errorf("location not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(location, id, "location"); err != nil {
		return fmt.Errorf("failed to delete location: %v", err)
	}

	return nil
}

func (m *LocationModel) GetAllLocations() ([]Location, error) {
	var locations []Location

	result, err := m.Service.GetAllEntities(&locations)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %v", err)
	}

	locationsPtr, ok := result.(*[]Location)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *locationsPtr, nil
}

func (m *LocationModel) GetLocationSubLocationsByFieldPreloaded(field, value string) (*Location, error) {
	var location Location

	preloads := []string{"SubLocations"}

	result, err := m.Service.GetEntityByFieldWithPreload(&location, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching location with sublocations by %s: %v", field, err)
		return nil, err
	}

	locationPtr, ok := result.(*Location)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return locationPtr, nil
}

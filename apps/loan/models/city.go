package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type City struct {
	ID uint `gorm:"primaryKey"`

	CityName string `gorm:"unique;not null;index"`
	RegionID uint   `gorm:"index"`
	Region   Region `gorm:"foreignKey:RegionID"`
}

type CityModel struct {
	Service services.Service
}

func NewCityModel(service services.Service) *CityModel {
	return &CityModel{Service: service}
}

// Uses the services interface which in turn uses the repository interface
func (m *CityModel) CreateCity(city *City) error {
	if err := m.Service.CreateEntity(city); err != nil {
		return fmt.Errorf("failed to create city: %v", err)
	}

	return nil
}

func (m *CityModel) DeleteCity(id uint) error {
	city := &City{}

	_, err := m.Service.GetEntityByID(city, id)
	if err != nil {
		return fmt.Errorf("city not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(city, id, "city"); err != nil {
		return fmt.Errorf("failed to delete city: %v", err)
	}

	return nil
}

func (m *CityModel) GetCityByField(field, value string) (*City, error) {
	var city City

	result, err := m.Service.GetEntityByField(field, value, &city)
	if err != nil {
		log.Printf("Error fetching city by %s: %v", field, err)
		return nil, err
	}

	return result.(*City), nil
}

func (m *CityModel) GetCities(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]City, int64, int64, error) {
	searchColumns := []string{"city_name"}

	preloads := []string{}

	citiesResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&City{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get cities: %v", err)
	}

	var cities []City
	for _, city := range citiesResult {
		if c, ok := city.(*City); ok {
			cities = append(cities, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", city)
		}
	}

	return cities, totalCount, filteredCount, nil
}

func (m *CityModel) UpdateCity(id int, cityName string, regionId uint) (City, error) {
	city := &City{ID: uint(id)}

	_, err := m.Service.GetEntityByID(city, uint(id))
	if err != nil {
		return City{}, fmt.Errorf("city not found: %v", err)
	}

	city.CityName = parameters.TrimWhitespace(cityName)
	city.RegionID = regionId

	if err := m.Service.UpdateEntity(city); err != nil {
		return City{}, fmt.Errorf("failed to update city: %v", err)
	}

	return *city, nil
}

func (m *CityModel) GetAllCities() ([]City, error) {
	var cities []City

	result, err := m.Service.GetAllEntities(&cities)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities: %v", err)
	}

	citiesPtr, ok := result.(*[]City)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *citiesPtr, nil
}

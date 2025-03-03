package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Country struct {
	ID          uint   `gorm:"primaryKey"`
	CountryName string `gorm:"unique;not null;index"`

	Regions []Region `gorm:"foreignKey:CountryID;constraint:onDelete:CASCADE"`
}

type CountryModel struct {
	Service services.Service
}

func NewCountryModel(service services.Service) *CountryModel {
	return &CountryModel{Service: service}
}

// Uses the services interface which in turn uses the repository interface
func (m *CountryModel) CreateCountry(country *Country) error {
	if err := m.Service.CreateEntity(country); err != nil {
		return fmt.Errorf("failed to create country: %v", err)
	}

	return nil
}

func (m *CountryModel) GetCountries(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Country, int64, int64, error) {
	searchColumns := []string{"country_name"}

	preloads := []string{}

	countriesResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Country{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get countries: %v", err)
	}

	var countries []Country
	for _, country := range countriesResult {
		if c, ok := country.(*Country); ok {
			countries = append(countries, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", country)
		}
	}

	return countries, totalCount, filteredCount, nil
}

func (m *CountryModel) GetCountryByField(field, value string) (*Country, error) {
	var country Country

	result, err := m.Service.GetEntityByField(field, value, &country)
	if err != nil {
		log.Printf("Error fetching country by %s: %v", field, err)
		return nil, err
	}

	return result.(*Country), nil
}

func (m *CountryModel) UpdateCountry(id int, countryName string) (Country, error) {
	country := &Country{ID: uint(id)}

	_, err := m.Service.GetEntityByID(country, uint(id))
	if err != nil {
		return Country{}, fmt.Errorf("country not found: %v", err)
	}

	country.CountryName = parameters.TrimWhitespace(countryName)

	if err := m.Service.UpdateEntity(country); err != nil {
		return Country{}, fmt.Errorf("failed to update country: %v", err)
	}

	return *country, nil
}

func (m *CountryModel) DeleteCountry(id uint) error {
	country := &Country{}

	_, err := m.Service.GetEntityByID(country, id)
	if err != nil {
		return fmt.Errorf("country not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(country, id, "country"); err != nil {
		return fmt.Errorf("failed to delete country: %v", err)
	}

	return nil
}

func (m *CountryModel) GetAllCountries() ([]Country, error) {
	var countries []Country

	result, err := m.Service.GetAllEntities(&countries)
	if err != nil {
		return nil, fmt.Errorf("failed to get countries: %v", err)
	}

	countriesPtr, ok := result.(*[]Country)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *countriesPtr, nil
}

func (m *CountryModel) GetCountryRegionsByFieldPreloaded(field, value string) (*Country, error) {
	var country Country

	preloads := []string{"Regions"}

	result, err := m.Service.GetEntityByFieldWithPreload(&country, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching country with regions by %s: %v", field, err)
		return nil, err
	}

	countryPtr, ok := result.(*Country)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return countryPtr, nil
}

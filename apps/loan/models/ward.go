package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Ward struct {
	ID uint `gorm:"primaryKey"`

	WardName    string `gorm:"unique;not null;index"`
	SubCountyID uint   `gorm:"index"`

	SubCounty SubCounty  `gorm:"foreignKey:SubCountyID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	Locations []Location `gorm:"foreignKey:WardID;constraint:OnDelete:CASCADE"`
}

type WardModel struct {
	Service services.Service
}

func NewWardModel(service services.Service) *WardModel {
	return &WardModel{Service: service}
}

func (m *WardModel) CreateWard(ward *Ward) error {
	if err := m.Service.CreateEntity(ward); err != nil {
		return fmt.Errorf("failed to create ward: %v", err)
	}

	return nil
}

func (m *WardModel) GetWards(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Ward, int64, int64, error) {
	searchColumns := []string{"ward_name"}

	preloads := []string{}

	wardsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Ward{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get wards: %v", err)
	}

	var wards []Ward
	for _, ward := range wardsResult {
		if c, ok := ward.(*Ward); ok {
			wards = append(wards, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", ward)
		}
	}

	return wards, totalCount, filteredCount, nil
}

func (m *WardModel) GetWardByField(field, value string) (*Ward, error) {
	var ward Ward

	result, err := m.Service.GetEntityByField(field, value, &ward)
	if err != nil {
		log.Printf("Error fetching ward by %s: %v", field, err)
		return nil, err
	}

	return result.(*Ward), nil
}

func (m *WardModel) UpdateWard(id int, wardName string, subCountyId uint) (Ward, error) {
	ward := &Ward{ID: uint(id)}

	_, err := m.Service.GetEntityByID(ward, uint(id))
	if err != nil {
		return Ward{}, fmt.Errorf("ward not found: %v", err)
	}

	ward.WardName = parameters.TrimWhitespace(wardName)
	ward.SubCountyID = subCountyId

	if err := m.Service.UpdateEntity(ward); err != nil {
		return Ward{}, fmt.Errorf("failed to update ward: %v", err)
	}

	return *ward, nil
}

func (m *WardModel) DeleteWard(id uint) error {
	ward := &Ward{}

	_, err := m.Service.GetEntityByID(ward, id)
	if err != nil {
		return fmt.Errorf("ward not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(ward, id, "ward"); err != nil {
		return fmt.Errorf("failed to delete ward: %v", err)
	}

	return nil
}

func (m *WardModel) GetAllWards() ([]Ward, error) {
	var wards []Ward

	result, err := m.Service.GetAllEntities(&wards)
	if err != nil {
		return nil, fmt.Errorf("failed to get wards: %v", err)
	}

	wardsPtr, ok := result.(*[]Ward)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *wardsPtr, nil
}

func (m *WardModel) GetWardLocationsByFieldPreloaded(field, value string) (*Ward, error) {
	var ward Ward

	preloads := []string{"Locations"}

	result, err := m.Service.GetEntityByFieldWithPreload(&ward, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching ward with locations by %s: %v", field, err)
		return nil, err
	}

	wardPtr, ok := result.(*Ward)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return wardPtr, nil
}

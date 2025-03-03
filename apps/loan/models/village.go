package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Village struct {
	ID uint `gorm:"primaryKey"`

	VillageName   string `gorm:"unique;not null;index"`
	SubLocationID uint   `gorm:"index"`

	SubLocation SubLocation `gorm:"foreignKey:SubLocationID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	Roads       []Road      `gorm:"foreignKey:VillageID;constraint:OnDelete:CASCADE"`
}

type VillageModel struct {
	Service services.Service
}

func NewVillageModel(service services.Service) *VillageModel {
	return &VillageModel{Service: service}
}

func (m *VillageModel) CreateVillage(village *Village) error {
	if err := m.Service.CreateEntity(village); err != nil {
		return fmt.Errorf("failed to create village: %v", err)
	}

	return nil
}

func (m *VillageModel) GetVillages(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Village, int64, int64, error) {
	searchColumns := []string{"village_name"}

	preloads := []string{}

	villagesResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Village{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get villages: %v", err)
	}

	var villages []Village
	for _, village := range villagesResult {
		if c, ok := village.(*Village); ok {
			villages = append(villages, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", village)
		}
	}

	return villages, totalCount, filteredCount, nil
}

func (m *VillageModel) GetVillageByField(field, value string) (*Village, error) {
	var village Village

	result, err := m.Service.GetEntityByField(field, value, &village)
	if err != nil {
		log.Printf("Error fetching village by %s: %v", field, err)
		return nil, err
	}

	return result.(*Village), nil
}

func (m *VillageModel) UpdateVillage(id int, villageName string, subLocationId uint) (Village, error) {
	village := &Village{ID: uint(id)}

	_, err := m.Service.GetEntityByID(village, uint(id))
	if err != nil {
		return Village{}, fmt.Errorf("village not found: %v", err)
	}

	village.VillageName = parameters.TrimWhitespace(villageName)
	village.SubLocationID = subLocationId

	if err := m.Service.UpdateEntity(village); err != nil {
		return Village{}, fmt.Errorf("failed to update village: %v", err)
	}

	return *village, nil
}

func (m *VillageModel) DeleteVillage(id uint) error {
	village := &Village{}

	_, err := m.Service.GetEntityByID(village, id)
	if err != nil {
		return fmt.Errorf("village not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(village, id, "village"); err != nil {
		return fmt.Errorf("failed to delete village: %v", err)
	}

	return nil
}

func (m *VillageModel) GetAllVillages() ([]Village, error) {
	var villages []Village

	result, err := m.Service.GetAllEntities(&villages)
	if err != nil {
		return nil, fmt.Errorf("failed to get villages: %v", err)
	}

	villagesPtr, ok := result.(*[]Village)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *villagesPtr, nil
}

func (m *VillageModel) GetVillageRoadsByFieldPreloaded(field, value string) (*Village, error) {
	var village Village

	preloads := []string{"Roads"}

	result, err := m.Service.GetEntityByFieldWithPreload(&village, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching village with roads by %s: %v", field, err)
		return nil, err
	}

	villagePtr, ok := result.(*Village)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return villagePtr, nil
}

package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Unit struct {
	ID uint `gorm:"primaryKey"`

	UnitName string `gorm:"unique;not null;index"`
	PlotID   uint   `gorm:"index"`
	Plot     Plot   `gorm:"foreignKey:PlotID"`
}

type UnitModel struct {
	Service services.Service
}

func NewUnitModel(service services.Service) *UnitModel {
	return &UnitModel{Service: service}
}

func (m *UnitModel) CreateUnit(unit *Unit) error {
	if err := m.Service.CreateEntity(unit); err != nil {
		return fmt.Errorf("failed to create unit: %v", err)
	}

	return nil
}

func (m *UnitModel) DeleteUnit(id uint) error {
	unit := &Unit{}

	_, err := m.Service.GetEntityByID(unit, id)
	if err != nil {
		return fmt.Errorf("unit not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(unit, id, "unit"); err != nil {
		return fmt.Errorf("failed to delete unit: %v", err)
	}

	return nil
}

func (m *UnitModel) GetUnitByField(field, value string) (*Unit, error) {
	var unit Unit

	result, err := m.Service.GetEntityByField(field, value, &unit)
	if err != nil {
		log.Printf("Error fetching unit by %s: %v", field, err)
		return nil, err
	}

	return result.(*Unit), nil
}

func (m *UnitModel) GetUnits(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Unit, int64, int64, error) {
	searchColumns := []string{"unit_name"}

	preloads := []string{}

	unitsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Unit{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get units: %v", err)
	}

	var units []Unit
	for _, unit := range unitsResult {
		if c, ok := unit.(*Unit); ok {
			units = append(units, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", unit)
		}
	}

	return units, totalCount, filteredCount, nil
}

func (m *UnitModel) UpdateUnit(id int, unitName string, plotId uint) (Unit, error) {
	unit := &Unit{ID: uint(id)}

	_, err := m.Service.GetEntityByID(unit, uint(id))
	if err != nil {
		return Unit{}, fmt.Errorf("unit not found: %v", err)
	}

	unit.UnitName = parameters.TrimWhitespace(unitName)
	unit.PlotID = plotId

	if err := m.Service.UpdateEntity(unit); err != nil {
		return Unit{}, fmt.Errorf("failed to update unit: %v", err)
	}

	return *unit, nil
}

func (m *UnitModel) GetAllUnits() ([]Unit, error) {
	var units []Unit

	result, err := m.Service.GetAllEntities(&units)
	if err != nil {
		return nil, fmt.Errorf("failed to get units: %v", err)
	}

	unitsPtr, ok := result.(*[]Unit)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *unitsPtr, nil
}

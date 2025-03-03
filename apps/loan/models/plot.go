package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Plot struct {
	ID uint `gorm:"primaryKey"`

	PlotName string `gorm:"unique;not null;index"`
	RoadID   uint   `gorm:"index"`

	Road  Road   `gorm:"foreignKey:RoadID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	Units []Unit `gorm:"foreignKey:PlotID;constraint:OnDelete:CASCADE"`
}

type PlotModel struct {
	Service services.Service
}

func NewPlotModel(service services.Service) *PlotModel {
	return &PlotModel{Service: service}
}

func (m *PlotModel) CreatePlot(plot *Plot) error {
	if err := m.Service.CreateEntity(plot); err != nil {
		return fmt.Errorf("failed to create plot: %v", err)
	}

	return nil
}

func (m *PlotModel) GetPlots(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Plot, int64, int64, error) {
	searchColumns := []string{"plot_name"}

	preloads := []string{}

	plotsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Plot{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get plots: %v", err)
	}

	var plots []Plot
	for _, plot := range plotsResult {
		if c, ok := plot.(*Plot); ok {
			plots = append(plots, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", plot)
		}
	}

	return plots, totalCount, filteredCount, nil
}

func (m *PlotModel) GetPlotByField(field, value string) (*Plot, error) {
	var plot Plot

	result, err := m.Service.GetEntityByField(field, value, &plot)
	if err != nil {
		log.Printf("Error fetching plot by %s: %v", field, err)
		return nil, err
	}

	return result.(*Plot), nil
}

func (m *PlotModel) UpdatePlot(id int, plotName string, roadId uint) (Plot, error) {
	plot := &Plot{ID: uint(id)}

	_, err := m.Service.GetEntityByID(plot, uint(id))
	if err != nil {
		return Plot{}, fmt.Errorf("plot not found: %v", err)
	}

	plot.PlotName = parameters.TrimWhitespace(plotName)
	plot.RoadID = roadId

	if err := m.Service.UpdateEntity(plot); err != nil {
		return Plot{}, fmt.Errorf("failed to update plot: %v", err)
	}

	return *plot, nil
}

func (m *PlotModel) DeletePlot(id uint) error {
	plot := &Plot{}

	_, err := m.Service.GetEntityByID(plot, id)
	if err != nil {
		return fmt.Errorf("plot not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(plot, id, "plot"); err != nil {
		return fmt.Errorf("failed to delete plot: %v", err)
	}

	return nil
}

func (m *PlotModel) GetAllPlots() ([]Plot, error) {
	var plots []Plot

	result, err := m.Service.GetAllEntities(&plots)
	if err != nil {
		return nil, fmt.Errorf("failed to get plots: %v", err)
	}

	plotsPtr, ok := result.(*[]Plot)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *plotsPtr, nil
}

func (m *PlotModel) GetPlotUnitsByFieldPreloaded(field, value string) (*Plot, error) {
	var plot Plot

	preloads := []string{"Units"}

	result, err := m.Service.GetEntityByFieldWithPreload(&plot, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching plot with units by %s: %v", field, err)
		return nil, err
	}

	plotPtr, ok := result.(*Plot)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return plotPtr, nil
}

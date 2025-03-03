package models

import (
	"fmt"
	"log"

	"github.com/kifangamukundi/gm/libs/parameters"
	"github.com/kifangamukundi/gm/loan/services"
)

type Road struct {
	ID uint `gorm:"primaryKey"`

	RoadName  string `gorm:"unique;not null;index"`
	VillageID uint   `gorm:"index"`

	Village Village `gorm:"foreignKey:VillageID;constraint:onDelete:CASCADE;onUpdate:CASCADE"`
	Plots   []Plot  `gorm:"foreignKey:RoadID;constraint:OnDelete:CASCADE"`
}

type RoadModel struct {
	Service services.Service
}

func NewRoadModel(service services.Service) *RoadModel {
	return &RoadModel{Service: service}
}

func (m *RoadModel) CreateRoad(road *Road) error {
	if err := m.Service.CreateEntity(road); err != nil {
		return fmt.Errorf("failed to create road: %v", err)
	}

	return nil
}

func (m *RoadModel) GetRoads(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Road, int64, int64, error) {
	searchColumns := []string{"road_name"}

	preloads := []string{}

	roadsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Road{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get roads: %v", err)
	}

	var roads []Road
	for _, road := range roadsResult {
		if c, ok := road.(*Road); ok {
			roads = append(roads, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", road)
		}
	}

	return roads, totalCount, filteredCount, nil
}

func (m *RoadModel) GetRoadByField(field, value string) (*Road, error) {
	var road Road

	result, err := m.Service.GetEntityByField(field, value, &road)
	if err != nil {
		log.Printf("Error fetching road by %s: %v", field, err)
		return nil, err
	}

	return result.(*Road), nil
}

func (m *RoadModel) UpdateRoad(id int, roadName string, villageId uint) (Road, error) {
	road := &Road{ID: uint(id)}

	_, err := m.Service.GetEntityByID(road, uint(id))
	if err != nil {
		return Road{}, fmt.Errorf("road not found: %v", err)
	}

	road.RoadName = parameters.TrimWhitespace(roadName)
	road.VillageID = villageId

	if err := m.Service.UpdateEntity(road); err != nil {
		return Road{}, fmt.Errorf("failed to update road: %v", err)
	}

	return *road, nil
}

func (m *RoadModel) DeleteRoad(id uint) error {
	road := &Road{}

	_, err := m.Service.GetEntityByID(road, id)
	if err != nil {
		return fmt.Errorf("road not found: %v", err)
	}

	if err := m.Service.HardDeleteEntity(road, id, "road"); err != nil {
		return fmt.Errorf("failed to delete road: %v", err)
	}

	return nil
}

func (m *RoadModel) GetAllRoads() ([]Road, error) {
	var roads []Road

	result, err := m.Service.GetAllEntities(&roads)
	if err != nil {
		return nil, fmt.Errorf("failed to get roads: %v", err)
	}

	roadsPtr, ok := result.(*[]Road)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *roadsPtr, nil
}

func (m *RoadModel) GetRoadPlotsByFieldPreloaded(field, value string) (*Road, error) {
	var road Road

	preloads := []string{"Plots"}

	result, err := m.Service.GetEntityByFieldWithPreload(&road, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching road with plots by %s: %v", field, err)
		return nil, err
	}

	roadPtr, ok := result.(*Road)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return roadPtr, nil
}

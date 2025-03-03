package models

import (
	"fmt"
	"log"
	"time"

	"github.com/kifangamukundi/gm/loan/services"
)

type Agent struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"unique; index"`
	User   User `gorm:"foreignKey:UserID"`

	IsActive  bool       `gorm:"default:false"`
	LastLogin *time.Time `gorm:"default:null"`

	CountryID uint    `gorm:"index"`
	Country   Country `gorm:"foreignKey:CountryID"`

	RegionID uint   `gorm:"index"`
	Region   Region `gorm:"foreignKey:RegionID"`

	CityID uint `gorm:"index"`
	City   City `gorm:"foreignKey:CityID"`

	Groups  []Group  `gorm:"foreignKey:AgentID"`
	Members []Member `gorm:"foreignKey:AgentID"`
}

type AgentModel struct {
	Service services.Service
}

func NewAgentModel(service services.Service) *AgentModel {
	return &AgentModel{Service: service}
}

func (m *AgentModel) CreateAgent(agent *Agent) error {
	if err := m.Service.CreateEntity(agent); err != nil {
		return fmt.Errorf("failed to create agent: %v", err)
	}

	return nil
}

func (m *AgentModel) CountAgents(conditions map[string]interface{}) (int64, error) {
	count, err := m.Service.CountEntities(&Agent{}, conditions)
	if err != nil {
		log.Printf("Error counting agents: %v", err)
		return 0, fmt.Errorf("failed to count agents: %v", err)
	}
	return count, nil
}

func (m *AgentModel) GetAgents(skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) ([]Agent, int64, int64, error) {
	searchColumns := []string{"user.first_name", "user.last_name"}

	preloads := []string{"User"}

	agentsResult, totalCount, filteredCount, err := m.Service.GetEntitiesFiltered(&Agent{}, skip, limit, sortOrder, sortByColumn, searchRegex, filterCriteria, searchColumns, preloads)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get agents: %v", err)
	}

	var agents []Agent
	for _, agent := range agentsResult {
		if c, ok := agent.(*Agent); ok {
			agents = append(agents, *c)
		} else {
			return nil, 0, 0, fmt.Errorf("unexpected type in result: %T", agent)
		}
	}

	return agents, totalCount, filteredCount, nil
}

func (m *AgentModel) GetAllAgents() ([]Agent, error) {
	preloads := []string{"User"}

	var agents []Agent

	result, err := m.Service.GetAllEntitiesWithPreload(&agents, preloads...)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %v", err)
	}

	agentsPtr, ok := result.(*[]Agent)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return *agentsPtr, nil
}

func (m *AgentModel) GetAgentByField(field, value string) (*Agent, error) {
	var agent Agent

	result, err := m.Service.GetEntityByField(field, value, &agent)
	if err != nil {
		log.Printf("Error fetching agent by %s: %v", field, err)
		return nil, err
	}

	return result.(*Agent), nil
}

func (m *AgentModel) GetAgentByFieldPreloaded(field, value string) (*Agent, error) {
	var agent Agent

	preloads := []string{"User"}

	result, err := m.Service.GetEntityByFieldWithPreload(&agent, field, value, preloads...)
	if err != nil {
		log.Printf("Error fetching agent with user info by %s: %v", field, err)
		return nil, err
	}

	agentPtr, ok := result.(*Agent)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return agentPtr, nil
}

func (m *AgentModel) UpdateAgent(id int, isActive bool) (Agent, error) {
	var agent Agent
	result, err := m.Service.GetEntityByID(&agent, uint(id))
	if err != nil {
		return Agent{}, fmt.Errorf("agent not found: %v", err)
	}

	fetchedAgent, ok := result.(*Agent)
	if !ok {
		return Agent{}, fmt.Errorf("unexpected type for agent result: %T", result)
	}

	fetchedAgent.IsActive = isActive

	if err := m.Service.UpdateEntity(fetchedAgent); err != nil {
		return Agent{}, fmt.Errorf("failed to update agent: %v", err)
	}

	return *fetchedAgent, nil
}

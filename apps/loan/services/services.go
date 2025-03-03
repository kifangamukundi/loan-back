package services

import (
	"fmt"

	"github.com/kifangamukundi/gm/loan/loanrepository"
)

type Service interface {
	CreateEntity(entity interface{}) error
	GetMonthlyEntityCounts(model interface{}) ([]map[string]interface{}, error)
	HardDeleteEntity(entity interface{}, id uint, entityName string, associations ...string) error
	SoftDeleteEntity(entity interface{}, id uint, entityName string, associations ...string) error
	CountEntities(model interface{}, conditions map[string]interface{}) (int64, error)
	UpdateEntity(entity interface{}) error
	GetEntityByID(entity interface{}, id uint) (interface{}, error)
	GetEntityByField(field, value string, model interface{}) (interface{}, error)
	GetEntitiesFiltered(model interface{}, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error)
	GetEntitiesFilteredLoans(model interface{}, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error)
	GetEntitiesFilteredTest(model interface{}, agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error)
	GetEntitiesFilteredAgentMemberLoans(model interface{}, agentId, groupId, memberId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error)
	GetEntitiesFilteredTest2(model interface{}, groupId, agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error)
	GetEntitiesFilteredByField(model interface{}, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string, slugColumn, slugValue string, joinTables []string, joinConditions []string) ([]interface{}, int64, int64, error)
	GetAllEntities(model interface{}) (interface{}, error)
	GetAllEntitiesWithPreload(model interface{}, preload ...string) (interface{}, error)
	GetEntityByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error)
	GetAllEntititiesByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error)
	GetEntitiesByFields(model interface{}, fieldValues map[string]interface{}) (interface{}, error)
	EntityClearAssociation(model interface{}, association string) error
}

type EntityServiceImpl struct {
	Repository loanrepository.LoanRepositoryInterface
}

func NewEntityService(repo loanrepository.LoanRepositoryInterface) Service {
	return &EntityServiceImpl{Repository: repo}
}

func (s *EntityServiceImpl) GetMonthlyEntityCounts(model interface{}) ([]map[string]interface{}, error) {
	results, err := s.Repository.GetMonthlyCounts(model)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch monthly counts: %v", err)
	}
	return results, nil
}

func (s *EntityServiceImpl) CountEntities(model interface{}, conditions map[string]interface{}) (int64, error) {
	count, err := s.Repository.CountConditional(model, conditions)
	if err != nil {
		return 0, fmt.Errorf("failed to count entities: %v", err)
	}
	return count, nil
}

func (s *EntityServiceImpl) GetEntityByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error) {
	result, err := s.Repository.GetByFieldWithPreload(model, field, value, preload...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entity with field %s and value %s: %v", field, value, err)
	}
	return result, nil
}

func (s *EntityServiceImpl) GetAllEntititiesByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error) {
	result, err := s.Repository.GetAllByFieldWithPreload(model, field, value, preload...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entity with field %s and value %s: %v", field, value, err)
	}
	return result, nil
}

func (s *EntityServiceImpl) EntityClearAssociation(model interface{}, association string) error {
	if err := s.Repository.ClearAssociations(model, association); err != nil {
		return fmt.Errorf("error clearing association: %v", err)
	}
	return nil
}

func (s *EntityServiceImpl) GetEntitiesFiltered(model interface{}, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error) {
	results, totalCount, filteredCount, err := s.Repository.GetAllFiltered(skip, limit, sortOrder, sortByColumn, searchRegex, searchColumns, filterCriteria, model, preload)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get entities: %v", err)
	}

	return results, totalCount, filteredCount, nil
}

func (s *EntityServiceImpl) GetEntitiesFilteredLoans(model interface{}, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error) {
	results, totalCount, filteredCount, err := s.Repository.GetAllFilteredLoans(skip, limit, sortOrder, sortByColumn, searchRegex, searchColumns, filterCriteria, model, preload)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get entities: %v", err)
	}

	return results, totalCount, filteredCount, nil
}

func (s *EntityServiceImpl) GetEntitiesFilteredTest(model interface{}, agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error) {
	results, totalCount, filteredCount, err := s.Repository.GetAllFilteredTest(agentId, skip, limit, sortOrder, sortByColumn, searchRegex, searchColumns, filterCriteria, model, preload)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get entities: %v", err)
	}

	return results, totalCount, filteredCount, nil
}

func (s *EntityServiceImpl) GetEntitiesFilteredAgentMemberLoans(model interface{}, agentId, groupId, memberId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error) {
	results, totalCount, filteredCount, err := s.Repository.GetAllFilteredAgentMemberLoans(agentId, groupId, memberId, skip, limit, sortOrder, sortByColumn, searchRegex, searchColumns, filterCriteria, model, preload)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get entities: %v", err)
	}

	return results, totalCount, filteredCount, nil
}

func (s *EntityServiceImpl) GetEntitiesFilteredTest2(model interface{}, groupId, agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string) ([]interface{}, int64, int64, error) {
	results, totalCount, filteredCount, err := s.Repository.GetAllFilteredTest2(groupId, agentId, skip, limit, sortOrder, sortByColumn, searchRegex, searchColumns, filterCriteria, model, preload)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get entities: %v", err)
	}

	return results, totalCount, filteredCount, nil
}

func (s *EntityServiceImpl) GetEntitiesFilteredByField(model interface{}, skip, limit int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}, searchColumns []string, preload []string, slugColumn, slugValue string, joinTables []string, joinConditions []string) ([]interface{}, int64, int64, error) {
	results, totalCount, filteredCount, err := s.Repository.GetAllFilteredByField(skip, limit, sortOrder, sortByColumn, searchRegex, searchColumns, filterCriteria, model, preload, slugColumn, slugValue, joinTables, joinConditions)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get entities: %v", err)
	}

	return results, totalCount, filteredCount, nil
}

func (s *EntityServiceImpl) GetAllEntities(model interface{}) (interface{}, error) {
	result, err := s.Repository.GetAll(model)
	if err != nil {
		return nil, fmt.Errorf("failed to get all entities: %v", err)
	}
	return result, nil
}

// test func
func (s *EntityServiceImpl) GetAllEntitiesWithPreload(model interface{}, preload ...string) (interface{}, error) {
	result, err := s.Repository.GetAllWithPreload(model, preload...)
	if err != nil {
		return nil, fmt.Errorf("failed to get all entities: %v", err)
	}
	return result, nil
}

func (s *EntityServiceImpl) CreateEntity(entity interface{}) error {
	if err := s.Repository.Create(entity); err != nil {
		return fmt.Errorf("failed to create entity: %v", err)
	}
	return nil
}

func (s *EntityServiceImpl) HardDeleteEntity(entity interface{}, id uint, entityName string, associations ...string) error {
	for _, assoc := range associations {
		if err := s.Repository.ClearAssociations(entity, assoc); err != nil {
			return fmt.Errorf("failed to clear association %s: %v", assoc, err)
		}
	}

	if err := s.Repository.HardDelete(entity, id); err != nil {
		return fmt.Errorf("failed to delete %s: %v", entityName, err)
	}

	return nil
}

func (s *EntityServiceImpl) SoftDeleteEntity(entity interface{}, id uint, entityName string, associations ...string) error {

	for _, assoc := range associations {
		if err := s.Repository.ClearAssociations(entity, assoc); err != nil {
			return fmt.Errorf("failed to clear association %s: %v", assoc, err)
		}
	}

	if err := s.Repository.SoftDelete(entity, id); err != nil {
		return fmt.Errorf("failed to delete %s: %v", entityName, err)
	}

	return nil
}

func (s *EntityServiceImpl) GetEntityByID(entity interface{}, id uint) (interface{}, error) {
	result, err := s.Repository.GetByID(entity, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity by ID: %v", err)
	}
	return result, nil
}

func (s *EntityServiceImpl) UpdateEntity(entity interface{}) error {
	if err := s.Repository.Update(entity); err != nil {
		return fmt.Errorf("failed to update entity: %v", err)
	}
	return nil
}

func (s *EntityServiceImpl) GetEntityByField(field, value string, model interface{}) (interface{}, error) {
	result, err := s.Repository.GetByField(model, field, value)
	if err != nil {
		return nil, fmt.Errorf("failed to get model by %s: %v", field, err)
	}
	return result, nil
}

func (s *EntityServiceImpl) GetEntitiesByFields(model interface{}, fieldValues map[string]interface{}) (interface{}, error) {
	if len(fieldValues) == 0 {
		return nil, fmt.Errorf("field values cannot be empty")
	}

	result, err := s.Repository.GetByFields(model, fieldValues)
	if err != nil {
		return nil, fmt.Errorf("error fetching entities: %v", err)
	}

	return result, nil
}

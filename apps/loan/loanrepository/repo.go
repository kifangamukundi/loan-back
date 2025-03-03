package loanrepository

import (
	"fmt"
	"reflect"

	"github.com/kifangamukundi/gm/libs/queryparams"
	"github.com/kifangamukundi/gm/libs/repositories"
	"gorm.io/gorm"
)

type LoanRepositoryInterface interface {
	repositories.Repository
	GetLoansByStatus(status string) ([]interface{}, error)
	GetAllFilteredLoans(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error)
	GetAllFilteredAgentMemberLoans(agentId, groupId, memberId, skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error)
	GetAllFilteredTest(agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error)
	GetAllFilteredTest2(groupId, agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error)
}

type LoanRepository struct {
	repo repositories.Repository
	DB   *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepositoryInterface {
	return &LoanRepository{
		repo: repositories.NewGormRepository(db),
		DB:   db,
	}
}

// ✅ Implement Custom Method
func (r *LoanRepository) GetLoansByStatus(status string) ([]interface{}, error) {
	var loans []interface{}
	err := r.DB.Where("status = ?", status).Find(&loans).Error
	return loans, err
}

func (r *LoanRepository) GetAllFilteredLoans(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error) {
	var result []interface{}
	var query = r.DB.Model(model)

	for _, assoc := range preload {
		query = query.Preload(assoc)
	}

	query = query.Where("status IN (?)", []string{"pending"})

	if searchRegex != "" {
		searchCondition := queryparams.BuildSearchCondition(searchColumns, searchRegex)
		query = query.Where(searchCondition)
	}

	if filterCriteria != nil {
		for key, value := range filterCriteria.(map[string]interface{}) {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	totalCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	filteredCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	query = query.Order(fmt.Sprintf("%s %s", sortByColumn, sortOrder)).
		Offset(skip).
		Limit(limit)

	resultPtr := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
	if err := query.Find(resultPtr.Interface()).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching models: %v", err)
	}

	resultSlice := resultPtr.Elem()
	for i := 0; i < resultSlice.Len(); i++ {
		result = append(result, resultSlice.Index(i).Interface())
	}

	return result, totalCount, filteredCount, nil
}

func (r *LoanRepository) GetAllFilteredAgentMemberLoans(agentId, groupId, memberId, skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error) {
	var result []interface{}
	var query = r.DB.Model(model)

	for _, assoc := range preload {
		query = query.Preload(assoc)
	}

	if searchRegex != "" {
		searchCondition := queryparams.BuildSearchCondition(searchColumns, searchRegex)
		query = query.Where(searchCondition)
	}

	if filterCriteria != nil {
		for key, value := range filterCriteria.(map[string]interface{}) {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	if agentId > 0 {
		query = query.Where("agent_id = ?", agentId)
	}

	if groupId > 0 {
		query = query.Where("group_id = ?", groupId)
	}

	if memberId > 0 {
		query = query.Where("member_id = ?", memberId)
	}

	totalCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	filteredCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	query = query.Order(fmt.Sprintf("%s %s", sortByColumn, sortOrder)).
		Offset(skip).
		Limit(limit)

	resultPtr := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
	if err := query.Find(resultPtr.Interface()).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching models: %v", err)
	}

	resultSlice := resultPtr.Elem()
	for i := 0; i < resultSlice.Len(); i++ {
		result = append(result, resultSlice.Index(i).Interface())
	}

	return result, totalCount, filteredCount, nil
}

func (r *LoanRepository) GetAllFilteredTest(agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error) {
	var result []interface{}
	var query = r.DB.Model(model)

	for _, assoc := range preload {
		query = query.Preload(assoc)
	}

	if searchRegex != "" {
		searchCondition := queryparams.BuildSearchCondition(searchColumns, searchRegex)
		query = query.Where(searchCondition)
	}

	if filterCriteria != nil {
		for key, value := range filterCriteria.(map[string]interface{}) {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	if agentId > 0 {
		query = query.Where("agent_id = ?", agentId)
	}

	totalCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	filteredCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	query = query.Order(fmt.Sprintf("%s %s", sortByColumn, sortOrder)).
		Offset(skip).
		Limit(limit)

	resultPtr := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
	if err := query.Find(resultPtr.Interface()).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching models: %v", err)
	}

	resultSlice := resultPtr.Elem()
	for i := 0; i < resultSlice.Len(); i++ {
		result = append(result, resultSlice.Index(i).Interface())
	}

	return result, totalCount, filteredCount, nil
}

func (r *LoanRepository) GetAllFilteredTest2(groupId, agentId, skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error) {
	var result []interface{}
	var query = r.DB.Model(model)

	for _, assoc := range preload {
		query = query.Preload(assoc)
	}

	if searchRegex != "" {
		searchCondition := queryparams.BuildSearchCondition(searchColumns, searchRegex)
		query = query.Where(searchCondition)
	}

	if filterCriteria != nil {
		for key, value := range filterCriteria.(map[string]interface{}) {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	// Filter by Group (Many-to-Many relationship with Groups)
	if groupId > 0 {
		// Join with group_members to filter by group_id
		query = query.Joins("JOIN group_members ON group_members.member_id = members.id").
			Where("group_members.group_id = ?", groupId)
	}

	if agentId > 0 {
		query = query.Where("agent_id = ?", agentId)
	}

	totalCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	filteredCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	query = query.Order(fmt.Sprintf("%s %s", sortByColumn, sortOrder)).
		Offset(skip).
		Limit(limit)

	resultPtr := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
	if err := query.Find(resultPtr.Interface()).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching models: %v", err)
	}

	resultSlice := resultPtr.Elem()
	for i := 0; i < resultSlice.Len(); i++ {
		result = append(result, resultSlice.Index(i).Interface())
	}

	return result, totalCount, filteredCount, nil
}

// ✅ Forward Base Repository Methods
func (r *LoanRepository) Create(model interface{}) error {
	return r.repo.Create(model)
}

func (r *LoanRepository) GetByID(model interface{}, id interface{}) (interface{}, error) {
	return r.repo.GetByID(model, id)
}

func (r *LoanRepository) Update(model interface{}) error {
	return r.repo.Update(model)
}

func (r *LoanRepository) SoftDelete(model interface{}, id interface{}) error {
	return r.repo.SoftDelete(model, id)
}

func (r *LoanRepository) HardDelete(model interface{}, id interface{}) error {
	return r.repo.HardDelete(model, id)
}

func (r *LoanRepository) GetByField(model interface{}, field string, value interface{}) (interface{}, error) {
	return r.repo.GetByField(model, field, value)
}

func (r *LoanRepository) GetByFields(model interface{}, fieldValues map[string]interface{}) (interface{}, error) {
	return r.repo.GetByFields(model, fieldValues)
}

func (r *LoanRepository) GetByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error) {
	return r.repo.GetByFieldWithPreload(model, field, value, preload...)
}

func (r *LoanRepository) GetAllByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error) {
	return r.repo.GetAllByFieldWithPreload(model, field, value, preload...)
}

func (r *LoanRepository) ClearAssociations(model interface{}, association string) error {
	return r.repo.ClearAssociations(model, association)
}

func (r *LoanRepository) GetAll(model interface{}) (interface{}, error) {
	return r.repo.GetAll(model)
}

func (r *LoanRepository) GetAllWithPreload(model interface{}, preload ...string) (interface{}, error) {
	return r.repo.GetAllWithPreload(model, preload...)
}

func (r *LoanRepository) GetAllFiltered(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error) {
	return r.repo.GetAllFiltered(skip, limit, sortOrder, sortByColumn, searchRegex, searchColumns, filterCriteria, model, preload)
}

func (r *LoanRepository) GetAllFilteredByField(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string, slugColumn, slugValue string, joinTables []string, joinConditions []string) ([]interface{}, int64, int64, error) {
	return r.repo.GetAllFilteredByField(skip, limit, sortOrder, sortByColumn, searchRegex, searchColumns, filterCriteria, model, preload, slugColumn, slugValue, joinTables, joinConditions)
}

func (r *LoanRepository) Count(query *gorm.DB, model interface{}) (int64, error) {
	return r.repo.Count(query, model)
}

func (r *LoanRepository) CountConditional(model interface{}, conditions map[string]interface{}) (int64, error) {
	return r.repo.CountConditional(model, conditions)
}

func (r *LoanRepository) Find(query *gorm.DB, model interface{}) ([]interface{}, error) {
	return r.repo.Find(query, model)
}

func (r *LoanRepository) GetMonthlyCounts(model interface{}) ([]map[string]interface{}, error) {
	return r.repo.GetMonthlyCounts(model)
}

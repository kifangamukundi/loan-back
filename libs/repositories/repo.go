package repositories

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/kifangamukundi/gm/libs/queryparams"

	"gorm.io/gorm"
)

type Repository interface {
	Create(model interface{}) error
	GetByID(model interface{}, id interface{}) (interface{}, error)
	GetMonthlyCounts(model interface{}) ([]map[string]interface{}, error)
	Update(model interface{}) error
	SoftDelete(model interface{}, id interface{}) error
	HardDelete(model interface{}, id interface{}) error
	GetByField(model interface{}, field string, value interface{}) (interface{}, error)
	GetByFields(model interface{}, fieldValues map[string]interface{}) (interface{}, error)
	GetByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error)
	GetAllByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error)
	ClearAssociations(model interface{}, association string) error
	GetAll(model interface{}) (interface{}, error)
	GetAllWithPreload(model interface{}, preload ...string) (interface{}, error)
	GetAllFiltered(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error)
	GetAllFilteredByField(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string, slugColumn, slugValue string, joinTables []string, joinConditions []string) ([]interface{}, int64, int64, error)
	Count(query *gorm.DB, model interface{}) (int64, error)
	CountConditional(model interface{}, conditions map[string]interface{}) (int64, error)
	Find(query *gorm.DB, model interface{}) ([]interface{}, error)
}

type GormRepository struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{DB: db}
}

func (r *GormRepository) GetMonthlyCounts(model interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.DB.Model(model).
		Select("TO_CHAR(created_at, 'Mon') AS month, COUNT(*) AS count, EXTRACT(MONTH FROM created_at) AS month_num").
		Group("TO_CHAR(created_at, 'Mon'), month_num").
		Order("month_num").
		Scan(&results).Error

	if err != nil {
		log.Printf("Error fetching monthly counts: %v", err)
		return nil, err
	}

	return results, nil
}

func (r *GormRepository) GetByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error) {
	query := r.DB.Where(field+" = ?", value)

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.First(model).Error
	if err != nil {
		log.Printf("Error fetching record by %s: %v", field, err)
		return nil, err
	}
	return model, nil
}

func (r *GormRepository) GetAllWithPreload(model interface{}, preload ...string) (interface{}, error) {
	query := r.DB

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Find(model).Error
	if err != nil {
		log.Printf("Error fetching records: %v", err)
		return nil, err
	}
	return model, nil
}

func (r *GormRepository) GetAllByFieldWithPreload(model interface{}, field, value string, preload ...string) (interface{}, error) {
	query := r.DB.Where(field+" = ?", value)

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Find(model).Error
	if err != nil {
		log.Printf("Error fetching record by %s: %v", field, err)
		return nil, err
	}
	return model, nil
}

func (r *GormRepository) Create(model interface{}) error {
	if err := r.DB.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create model: %v", err)
	}
	return nil
}

func (r *GormRepository) GetByID(model interface{}, id interface{}) (interface{}, error) {
	if err := r.DB.Where("id = ?", id).First(model).Error; err != nil {
		return nil, fmt.Errorf("failed to get model by ID: %v", err)
	}
	return model, nil
}

func (r *GormRepository) Update(model interface{}) error {
	if err := r.DB.Save(model).Error; err != nil {
		return fmt.Errorf("failed to update model: %v", err)
	}
	return nil
}

func (r *GormRepository) SoftDelete(model interface{}, id interface{}) error {
	if err := r.DB.Model(model).Where("id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to soft delete model: %v", err)
	}
	return nil
}

func (r *GormRepository) HardDelete(model interface{}, id interface{}) error {
	if err := r.DB.Unscoped().Delete(model, id).Error; err != nil {
		return fmt.Errorf("failed to hard delete model: %v", err)
	}
	return nil
}

func (r *GormRepository) GetByField(model interface{}, field string, value interface{}) (interface{}, error) {
	if err := r.DB.Where(fmt.Sprintf("%s = ?", field), value).First(model).Error; err != nil {
		return nil, fmt.Errorf("failed to get model by %s: %v", field, err)
	}
	return model, nil
}

func (r *GormRepository) GetByFields(model interface{}, fieldValues map[string]interface{}) (interface{}, error) {
	if len(fieldValues) == 0 {
		return nil, fmt.Errorf("field values cannot be empty")
	}

	query := r.DB
	for field, value := range fieldValues {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			query = query.Where(fmt.Sprintf("%s IN (?)", field), value)
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	if err := query.Find(model).Error; err != nil {
		return nil, fmt.Errorf("failed to get model by fields: %v", err)
	}
	return model, nil
}

func (r *GormRepository) GetAll(model interface{}) (interface{}, error) {
	if err := r.DB.Find(model).Error; err != nil {
		return nil, fmt.Errorf("failed to get all items: %v", err)
	}
	return model, nil
}

func (r *GormRepository) ClearAssociations(model interface{}, association string) error {
	if err := r.DB.Model(model).Association(association).Clear(); err != nil {
		return fmt.Errorf("failed to clear association %s: %v", association, err)
	}
	return nil
}

func (r *GormRepository) Count(query *gorm.DB, model interface{}) (int64, error) {
	var count int64
	if err := query.Model(model).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting models: %v", err)
	}
	return count, nil
}

func (r *GormRepository) CountConditional(model interface{}, conditions map[string]interface{}) (int64, error) {
	var count int64
	query := r.DB.Model(model)

	// Apply conditions if provided
	if len(conditions) > 0 {
		query = query.Where(conditions)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting models: %v", err)
	}
	return count, nil
}

func (r *GormRepository) Find(query *gorm.DB, model interface{}) ([]interface{}, error) {
	var result []interface{}
	if err := query.Find(&result).Error; err != nil {
		return nil, fmt.Errorf("error fetching models: %v", err)
	}
	return result, nil
}

func (r *GormRepository) GetAllFiltered(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string) ([]interface{}, int64, int64, error) {
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

func (r *GormRepository) GetAllFilteredByField(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria interface{}, model interface{}, preload []string, slugColumn, slugValue string, joinTables []string, joinConditions []string) ([]interface{}, int64, int64, error) {
	var result []interface{}
	var query = r.DB.Model(model)

	for _, assoc := range preload {
		query = query.Preload(assoc)
	}

	if slugColumn != "" && slugValue != "" {
		for i, table := range joinTables {
			query = query.Joins(fmt.Sprintf("JOIN %s ON %s", table, joinConditions[i]))
		}
		query = query.Where(fmt.Sprintf("%s = ?", slugColumn), slugValue)
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

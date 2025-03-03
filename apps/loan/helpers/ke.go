package helpers

import "fmt"

func GenerateCacheKey(itemName string, page, limit, skip int, sortOrder, sortByColumn, searchRegex string, filterCriteria interface{}) string {
	filterCriteriaStr := fmt.Sprintf("%v", filterCriteria)

	return fmt.Sprintf("%s:%d:%d:%d:%s:%s:%s:%s", itemName, page, limit, skip, sortOrder, sortByColumn, searchRegex, filterCriteriaStr)
}

func GenerateUserPermissionsCacheKey(itemName string, id uint) string {
	return fmt.Sprintf("%s:%v", itemName, id)
}

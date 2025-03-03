package queryparams

import "github.com/gin-gonic/gin"

func SortColumnMiddleware(validSortCriteria []string, defaultSortCriteria string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sortByColumn := c.DefaultQuery("sortByColumn", defaultSortCriteria)
		if !contains(validSortCriteria, sortByColumn) {
			sortByColumn = defaultSortCriteria
		}
		c.Set("sortByColumn", sortByColumn)
		c.Next()
	}
}

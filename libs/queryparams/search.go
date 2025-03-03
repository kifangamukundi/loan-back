package queryparams

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func SearchMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.DefaultQuery("search", "")
		if len(search) >= 3 {
			c.Set("searchRegex", "%"+search+"%")
		} else {
			c.Set("searchRegex", "")
		}
		c.Next()
	}
}

func BuildSearchCondition(columns []string, searchRegex string) string {
	searchConditions := make([]string, len(columns))
	for i, column := range columns {
		searchConditions[i] = fmt.Sprintf("%s ILIKE '%%%s%%'", column, searchRegex)
	}
	return fmt.Sprintf("(%s)", strings.Join(searchConditions, " OR "))
}

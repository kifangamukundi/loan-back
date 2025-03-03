package queryparams

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func SortOrderMiddleware(validSortOrders []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sortOrder := strings.ToLower(c.DefaultQuery("sortOrder", "asc"))
		if !contains(validSortOrders, sortOrder) {
			sortOrder = "asc"
		}
		c.Set("sortOrder", sortOrder)
		c.Next()
	}
}

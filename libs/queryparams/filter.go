package queryparams

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		filters := c.DefaultQuery("filters", "{}")
		var filterCriteria map[string]interface{}
		if err := json.Unmarshal([]byte(filters), &filterCriteria); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter criteria"})
			c.Abort()
			return
		}
		c.Set("filterCriteria", filterCriteria)
		c.Next()
	}
}

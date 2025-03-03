package queryparams

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func PaginationMiddleware(defaultPage, defaultLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = defaultPage
		}
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil || limit < 1 {
			limit = defaultLimit
		}
		c.Set("pagination", gin.H{
			"limit": limit,
			"skip":  (page - 1) * limit,
			"page":  page,
		})
		c.Next()
	}
}

package binders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindJSON is a helper function to bind JSON request body
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return err
	}
	return nil
}

// BindXML is a helper function to bind XML request body
func BindXML(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindXML(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return err
	}
	return nil
}

// BindForm is a helper function to bind form data request body
func BindForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return err
	}
	return nil
}

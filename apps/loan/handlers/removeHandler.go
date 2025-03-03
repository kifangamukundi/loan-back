package handlers

import (
	"context"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func RemoveMedia(cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			PublicID string `json:"public_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Public ID is required"})
			return
		}

		// Create a context for the Cloudinary operation
		ctx := context.Background()

		// Remove media by public ID
		result, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: request.PublicID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete Success",
			"data":    result,
		})
	}
}

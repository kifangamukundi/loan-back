package routes

import (
	"github.com/kifangamukundi/gm/loan/config"
	"github.com/kifangamukundi/gm/loan/handlers"
	"github.com/kifangamukundi/gm/loan/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MediaRoutes(r *gin.Engine, db *gorm.DB) {
	cloudConfig, cld := config.GetCloudinaryConfig()
	
	api := r.Group("/api")

	v1 := api.Group("/v1/media")
	{
		v1.POST("/new", middlewares.AdvancedAuth(db, []string{"upload_media"}), handlers.UploadMedia(cloudConfig, cld))
		v1.POST("/remove", middlewares.AdvancedAuth(db, []string{"delete_media"}), handlers.RemoveMedia(cld))
	}

	v2 := api.Group("/v2/media")
	{
		v2.POST("/new", handlers.UploadMedia(cloudConfig, cld))
	}
}

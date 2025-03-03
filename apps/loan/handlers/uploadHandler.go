package handlers

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/kifangamukundi/gm/libs/binders"
	"github.com/kifangamukundi/gm/loan/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type UploadResponse struct {
	PublicID  string `json:"public_id"`
	SecureURL string `json:"secure_url"`
}

func UploadMedia(config *config.CloudinaryConfig, cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form-data"})
			return
		}

		files := form.File["files"]
		subfolder := c.DefaultPostForm("subfolder", "default_subfolder")
		folderName := filepath.Join(config.RootFolder, subfolder)

		var results []UploadResponse

		ctx := context.Background()

		for _, file := range files {
			uploadResult, err := uploadFile(ctx, cld, file, folderName)
			if err != nil {
				log.Println("Error uploading file:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload some files"})
				return
			}

			result := UploadResponse{
				PublicID:  uploadResult.PublicID,
				SecureURL: uploadResult.SecureURL,
			}
			results = append(results, result)
		}

		binders.ReturnJSONMediaUploadResponse(c, results)
	}
}

func uploadFile(ctx context.Context, cld *cloudinary.Cloudinary, file *multipart.FileHeader, folderName string) (*uploader.UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	uploadParams := uploader.UploadParams{Folder: folderName}
	return cld.Upload.Upload(ctx, src, uploadParams)
}

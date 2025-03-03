package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

type CloudinaryConfig struct {
	CloudName  string
	APIKey     string
	APISecret  string
	RootFolder string
}

func GetCloudinaryConfig() (*CloudinaryConfig, *cloudinary.Cloudinary) {
	config := &CloudinaryConfig{
		CloudName:  os.Getenv("CLOUDINARY_CLOUD_NAME"),
		APIKey:     os.Getenv("CLOUDINARY_API_KEY"),
		APISecret:  os.Getenv("CLOUDINARY_API_SECRET"),
		RootFolder: os.Getenv("CLOUDINARY_ROOT_FOLDER"),
	}

	cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	return config, cld
}

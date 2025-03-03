package handlers

import (
	"context"
	"fmt"

	"github.com/kifangamukundi/gm/loan/config"
	"github.com/kifangamukundi/gm/loan/deserializers"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// RemoveAllImages deletes a list of images from Cloudinary by their PublicIDs.
func RemoveAllImages(config *config.CloudinaryConfig, cld *cloudinary.Cloudinary, images []deserializers.DefaultImage) error {
	ctx := context.Background()

	// You can use config here if you need to refer to settings like RootFolder, etc.

	for _, img := range images {
		if img.PublicID == nil {
			return fmt.Errorf("image does not have a valid PublicID")
		}

		_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: *img.PublicID})
		if err != nil {
			return fmt.Errorf("failed to delete image (%s): %v", *img.PublicID, err)
		}
	}
	return nil
}

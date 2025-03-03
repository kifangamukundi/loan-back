package deserializers

import (
	"encoding/json"
	"errors"
	"fmt"
)

type DefaultImage struct {
	SecureUrl *string `json:"secure_url"`
	PublicID  *string `json:"public_id"`
}

// Define a type alias for []DefaultImage to implement custom methods
type DefaultImageSlice []DefaultImage

// Implement the Scanner interface for DefaultImageSlice
func (di *DefaultImageSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Try unmarshalling as an array
	if err := json.Unmarshal(bytes, di); err == nil {
		return nil
	}

	// If unmarshalling as an array fails, try unmarshalling as a single object
	var single DefaultImage
	if err := json.Unmarshal(bytes, &single); err != nil {
		return fmt.Errorf("failed to unmarshal DefaultImageSlice: %v", err)
	}

	// If single object, wrap it in a slice
	*di = DefaultImageSlice{single}
	return nil
}

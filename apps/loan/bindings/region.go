package bindings

type CreateRegionRequest struct {
	RegionName string `json:"RegionName" binding:"required,min=3,max=100"`
	CountryID  int    `json:"CountryID" binding:"required"`
}

type UpdateRegionRequest struct {
	RegionName string `json:"RegionName" binding:"required,min=3,max=100"`
	CountryID  int    `json:"CountryID" binding:"required"`
}

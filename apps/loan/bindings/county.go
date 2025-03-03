package bindings

type CreateCountyRequest struct {
	CountyName string `json:"CountyName" binding:"required,min=3,max=100"`
	RegionID  int    `json:"RegionID" binding:"required"`
}

type UpdateCountyRequest struct {
	CountyName string `json:"CountyName" binding:"required,min=3,max=100"`
	RegionID  int    `json:"RegionID" binding:"required"`
}
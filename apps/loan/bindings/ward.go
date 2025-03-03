package bindings

type CreateWardRequest struct {
	WardName string `json:"WardName" binding:"required,min=3,max=100"`
	SubCountyID  int    `json:"SubCountyID" binding:"required"`
}

type UpdateWardRequest struct {
	WardName string `json:"WardName" binding:"required,min=3,max=100"`
	SubCountyID  int    `json:"SubCountyID" binding:"required"`
}

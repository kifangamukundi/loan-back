package bindings

type CreateSubLocationRequest struct {
	SubLocationName string `json:"SubLocationName" binding:"required,min=3,max=100"`
	LocationID  int    `json:"LocationID" binding:"required"`
}

type UpdateSubLocationRequest struct {
	SubLocationName string `json:"SubLocationName" binding:"required,min=3,max=100"`
	LocationID  int    `json:"LocationID" binding:"required"`
}
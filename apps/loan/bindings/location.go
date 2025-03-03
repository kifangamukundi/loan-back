package bindings

type CreateLocationRequest struct {
	LocationName string `json:"LocationName" binding:"required,min=3,max=100"`
	WardID  int    `json:"WardID" binding:"required"`
}

type UpdateLocationRequest struct {
	LocationName string `json:"LocationName" binding:"required,min=3,max=100"`
	WardID  int    `json:"WardID" binding:"required"`
}

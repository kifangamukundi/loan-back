package bindings

type CreateVillageRequest struct {
	VillageName string `json:"VillageName" binding:"required,min=3,max=100"`
	SubLocationID  int    `json:"SubLocationID" binding:"required"`
}

type UpdateVillageRequest struct {
	VillageName string `json:"VillageName" binding:"required,min=3,max=100"`
	SubLocationID  int    `json:"SubLocationID" binding:"required"`
}
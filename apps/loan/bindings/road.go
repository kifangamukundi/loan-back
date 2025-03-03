package bindings

type CreateRoadRequest struct {
	RoadName string `json:"RoadName" binding:"required,min=3,max=100"`
	VillageID  int    `json:"VillageID" binding:"required"`
}

type UpdateRoadRequest struct {
	RoadName string `json:"RoadName" binding:"required,min=3,max=100"`
	VillageID  int    `json:"VillageID" binding:"required"`
}
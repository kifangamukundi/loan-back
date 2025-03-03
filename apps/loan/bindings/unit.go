package bindings

type CreateUnitRequest struct {
	UnitName string `json:"UnitName" binding:"required,min=3,max=100"`
	PlotID int    `json:"PlotID" binding:"required"`
}

type UpdateUnitRequest struct {
	UnitName string `json:"UnitName" binding:"required,min=3,max=100"`
	PlotID int    `json:"PlotID" binding:"required"`
}
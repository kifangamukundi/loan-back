package bindings

type CreatePlotRequest struct {
	PlotName string `json:"PlotName" binding:"required,min=3,max=100"`
	RoadID  int    `json:"RoadID" binding:"required"`
}

type UpdatePlotRequest struct {
	PlotName string `json:"PlotName" binding:"required,min=3,max=100"`
	RoadID  int    `json:"RoadID" binding:"required"`
}
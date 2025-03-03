package bindings

// Request structures for user routes (json, xml, form)

type CreateGroup struct {
	GroupName    string `json:"GroupName" binding:"required,min=3,max=100"`
	AgentID    int    `json:"AgentID" binding:"required"`
	CountryID    int    `json:"CountryID" binding:"required"`
	RegionID     int    `json:"RegionID" binding:"required"`
	CityID       int    `json:"CityID" binding:"required"`
}

type UpdateGroup struct {
	IsActive bool  `json:"IsActive"`
	AgentID    int    `json:"AgentID" binding:"required"`
}

type GroupResponse struct {
	ID        uint   `json:"id"`
	GroupName string `json:"GroupName"`
	AgentID        uint   `json:"AgentID"`
	IsActive  bool   `json:"IsActive"`
}
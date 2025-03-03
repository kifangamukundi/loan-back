package bindings

type CreateSubCountyRequest struct {
	SubCountyName string `json:"SubCountyName" binding:"required,min=3,max=100"`
	CountyID  int    `json:"CountyID" binding:"required"`
}

type UpdateSubCountyRequest struct {
	SubCountyName string `json:"SubCountyName" binding:"required,min=3,max=100"`
	CountyID  int    `json:"CountyID" binding:"required"`
}

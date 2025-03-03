package bindings

// Request structures for user routes (json, xml, form)

type CreateMember struct {
	FirstName    string `json:"FirstName" binding:"required,min=3,max=100"`
	LastName     string `json:"LastName" binding:"required,min=3,max=100"`
	Email        string `json:"Email" binding:"required,email"`
	MobileNumber string `json:"MobileNumber" binding:"required"`
	CountryID    int    `json:"CountryID" binding:"required"`
	RegionID     int    `json:"RegionID" binding:"required"`
	CityID       int    `json:"CityID" binding:"required"`

	Groups []int `json:"Groups" binding:"required,dive,number"`
}

type UpdateMember struct {
	IsActive bool  `json:"IsActive"`
	Groups   []int `json:"Groups" binding:"required,dive,number"`
}

type MemberResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
	Mobile    string `json:"MobileNumber"`
	IsActive  bool   `json:"IsActive"`
}

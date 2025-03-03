package bindings

type CreateCountryRequest struct {
	CountryName string `json:"CountryName" binding:"required,min=3,max=100"`
}

type UpdateCountryRequest struct {
	CountryName string `json:"CountryName" binding:"required,min=3,max=100"`
}

package bindings

// Request structures for user routes (json, xml, form)

type RegisterRequest struct {
	FirstName    string `json:"FirstName" binding:"required,min=3,max=100"`
	LastName     string `json:"LastName" binding:"required,min=3,max=100"`
	Email        string `json:"Email" binding:"required,email"`
	MobileNumber string `json:"MobileNumber" binding:"required"`
	Password     string `json:"Password" binding:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"Email" binding:"required"`
}

type ChangePasswordRequest struct {
	Password string `json:"Password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"Email" binding:"required,email"`
	Password string `json:"Password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"RefreshToken" binding:"required"`
	UserID       int    `json:"UserID" binding:"required"`
}

type UpdateUserRequest struct {
	IsActive bool  `json:"IsActive"`
	IsLocked bool  `json:"IsLocked"`
	Roles    []int `json:"Roles" binding:"required,dive,number"`
}

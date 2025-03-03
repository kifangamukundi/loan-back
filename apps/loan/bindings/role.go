package bindings

type CreateRoleRequest struct {
	RoleName    string `json:"RoleName" binding:"required,min=3,max=100"`
	Permissions []int  `json:"Permissions" binding:"required,dive,number"`
}

type UpdateRoleRequest struct {
	RoleName    string `json:"RoleName" binding:"required,min=3,max=100"`
	Permissions []int  `json:"Permissions" binding:"required,dive,number"`
}

package bindings

type CreatePermissionRequest struct {
	PermissionName string `json:"PermissionName" binding:"required,min=3,max=100"`
	Roles          []int  `json:"Roles" binding:"required,dive,number"`
}

type UpdatePermissionRequest struct {
	PermissionName string `json:"PermissionName" binding:"required,min=3,max=100"`
	Roles          []int  `json:"Roles" binding:"required,dive,number"`
}

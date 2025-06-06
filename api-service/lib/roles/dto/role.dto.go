package roles

type RoleDTO struct {
	ID         int    `json:"id"`
	RoleName   string `json:"role_name" binding:"required"`
	RoleStatus int    `json:"role_status" binding:"required"`
	RoleKey    string `json:"role_key" binding:"required"`
}

type RoleCreateDTO struct {
	RoleName   string `json:"role_name" binding:"required"`
	RoleStatus int    `json:"role_status" binding:"required"`
	RoleKey    string `json:"role_key" binding:"required"`
}

type RoleUpdateDTO struct {
	ID         int    `json:"id"`
	RoleName   string `json:"role_name" binding:"required"`
	RoleStatus int    `json:"role_status" binding:"required"`
	RoleKey    string `json:"role_key" binding:"required"`
}

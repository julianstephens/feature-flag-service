package rbac

import "time"

type RbacRoleDto struct {
	RbacRole
	Permissions []RbacPermission `json:"permissions"`
}

type RbacRole struct {
	ID           int       `json:"-"`
	PublicRoleID string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateRbacRoleRequest struct {
	PublicRoleID string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type RbacUserDto struct {
	RbacUser
	Roles []string `json:"roles"`
}

type RbacUser struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRbacUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RbacPermission struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRbacPermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RbacRolePermission struct {
	RoleName       string `json:"role_name"`
	PermissionName string `json:"permission_name"`
}

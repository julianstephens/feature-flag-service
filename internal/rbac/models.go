package rbac

import "time"

type RbacRoleDto struct {
	RbacRole
	Permissions []RbacPermission `json:"permissions"`
}

type RbacRole struct {
	ID           int       `json:"-" db:"id"`
	PublicRoleID string    `json:"id" db:"public_id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
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
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"-" db:"password"`
	Activated bool      `json:"activated" db:"activated"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRbacUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UpdateRbacUserRequest struct {
	CreateRbacUserRequest
}

type RbacPermission struct {
	ID          int       `json:"-" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateRbacPermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RbacRolePermission struct {
	RoleName       string `json:"role_name"`
	PermissionName string `json:"permission_name"`
}

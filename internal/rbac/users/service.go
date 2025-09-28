package users

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"

	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/rbac"
	"github.com/julianstephens/feature-flag-service/internal/storage"
	"github.com/julianstephens/go-utils/helpers"
	"github.com/julianstephens/go-utils/logger"
	"github.com/julianstephens/go-utils/security"
)

// TODO: standardize role schema
type RbacUserService struct {
	conf  *config.Config
	store *storage.PostgresStore
}

const (
	RBAC_USER_TABLE       = "rbac_users"
	RBAC_USER_ROLES_TABLE = "rbac_user_roles"

	JOIN_RBAC_USER_ROLES = `SELECT r.*
		FROM rbac_roles r
		JOIN rbac_user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1;`

	JOIN_RBAC_ROLE_PERMISSIONS = `SELECT p.*
		FROM rbac_permissions p
		JOIN rbac_role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1;`
)

type Service interface {
	CreateUser(ctx context.Context, email, name, password string, roles []string) (*rbac.RbacUserDto, error)
	UpdateUser(ctx context.Context, id, email, name string, roles []string) error
	GetUser(ctx context.Context, id string) (*rbac.RbacUserDto, error)
	GetUserByEmail(ctx context.Context, email string) (*rbac.RbacUserDto, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]*rbac.RbacUserDto, error)
	ListUserRoles(ctx context.Context, id string) ([]*rbac.RbacRoleDto, error)
}

func NewRbacUserService(conf *config.Config, store *storage.PostgresStore) Service {
	return &RbacUserService{
		conf:  conf,
		store: store,
	}
}

func (s *RbacUserService) CreateUser(ctx context.Context, email, name, password string, roles []string) (*rbac.RbacUserDto, error) {
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := rbac.RbacUser{
		ID:       "user-" + uuid.New().String(),
		Email:    email,
		Name:     name,
		Password: hashedPassword,
	}

	if err := s.store.Post(ctx, RBAC_USER_TABLE, helpers.StructToMap(newUser)); err != nil {
		logger.Errorf("Error creating user %s: %v", email, err)
		return nil, err
	}

	s.store.Post(ctx, RBAC_USER_ROLES_TABLE, map[string]any{
		"user_id": newUser.ID,
		"roles":   roles,
	})

	return &rbac.RbacUserDto{
		RbacUser: newUser,
	}, nil
}

func (s *RbacUserService) GetUser(ctx context.Context, id string) (*rbac.RbacUserDto, error) {
	var rbacUser rbac.RbacUser
	if err := s.store.Get(ctx, RBAC_USER_TABLE, &rbacUser, "id=$1", id); err != nil {
		logger.Errorf("Error fetching user by id %s: %v", id, err)
		if pgxscan.NotFound(err) {
			return nil, storage.ErrKeyNotFound
		}
		return nil, err
	}

	var roleNames []string
	rows, err := s.store.Query(ctx, JOIN_RBAC_USER_ROLES, rbacUser.ID)
	if err != nil {
		logger.Errorf("Error fetching roles for user %s: %v", rbacUser.ID, err)
		if pgxscan.NotFound(err) {
			return nil, storage.ErrKeyNotFound
		}
		return nil, err
	}
	defer rows.Close()

	rs := pgxscan.NewRowScanner(rows)
	for rows.Next() {
		var role rbac.RbacRole
		if err := rs.Scan(&role); err != nil {
			logger.Errorf("Error scanning role for user %s: %v", rbacUser.ID, err)
			return nil, err
		}
		roleNames = append(roleNames, role.Name)
	}

	return &rbac.RbacUserDto{
		RbacUser: rbacUser,
		Roles:    roleNames,
	}, nil
}

func (s *RbacUserService) GetUserByEmail(ctx context.Context, email string) (*rbac.RbacUserDto, error) {
	var rbacUser rbac.RbacUser
	if err := s.store.Get(ctx, RBAC_USER_TABLE, &rbacUser, "email=$1", email); err != nil {
		logger.Errorf("Error fetching user by email %s: %v", email, err)
		if pgxscan.NotFound(err) {
			return nil, storage.ErrKeyNotFound
		}
		return nil, err
	}

	var roles []string
	rows, err := s.store.Query(ctx, JOIN_RBAC_USER_ROLES, rbacUser.ID)
	if err != nil {
		logger.Errorf("Error fetching roles for user %s: %v", rbacUser.ID, err)
		if pgxscan.NotFound(err) {
			return nil, storage.ErrKeyNotFound
		}
		return nil, err
	}
	defer rows.Close()

	rs := pgxscan.NewRowScanner(rows)
	for rows.Next() {
		var role rbac.RbacRole
		if err := rs.Scan(&role); err != nil {
			logger.Errorf("Error scanning role for user %s: %v", rbacUser.ID, err)
			return nil, err
		}
		roles = append(roles, role.Name)
	}

	return &rbac.RbacUserDto{
		RbacUser: rbacUser,
		Roles:    roles,
	}, nil
}

func (s *RbacUserService) UpdateUser(ctx context.Context, id, email, password string, roles []string) error {
	// return s.store.Put(ctx, id, map[string]interface{}{
	// 	"username": username,
	// 	"email":    email,
	// 	"password": password,
	// 	"roles":    roles,
	// }, s.rbacUserTable)
	return nil
}

func (s *RbacUserService) DeleteUser(ctx context.Context, id string) error {
	// return s.store.Delete(ctx, id, s.rbacUserTable)
	return nil
}

func (s *RbacUserService) ListUsers(ctx context.Context) ([]*rbac.RbacUserDto, error) {
	var users []*rbac.RbacUser
	err := s.store.ListAll(ctx, RBAC_USER_TABLE, &users)
	if err != nil {
		return nil, err
	}

	var userDtos []*rbac.RbacUserDto
	for _, u := range users {
		if u == nil {
			continue
		}

		user := *u

		var roles []string
		rows, err := s.store.Query(ctx, JOIN_RBAC_USER_ROLES, user.ID)
		if err != nil {
			logger.Errorf("Error fetching roles for user %s: %v", user.ID, err)
			if pgxscan.NotFound(err) {
				return nil, storage.ErrKeyNotFound
			}
			return nil, err
		}
		defer rows.Close()

		rs := pgxscan.NewRowScanner(rows)
		for rows.Next() {
			var role rbac.RbacRole
			if err := rs.Scan(&role); err != nil {
				logger.Errorf("Error scanning role for user %s: %v", user.ID, err)
				return nil, err
			}
			roles = append(roles, role.Name)
		}
		userDtos = append(userDtos, &rbac.RbacUserDto{
			RbacUser: user,
			Roles:    roles,
		})
	}

	return userDtos, nil
}

func (s *RbacUserService) ListUserRoles(ctx context.Context, id string) ([]*rbac.RbacRoleDto, error) {
	var roles []*rbac.RbacRole
	rows, err := s.store.Query(ctx, JOIN_RBAC_USER_ROLES, id)
	if err != nil {
		return nil, err
	}

	rs := pgxscan.NewRowScanner(rows)

	for rows.Next() {
		var role rbac.RbacRole
		if err := rs.Scan(&role); err != nil {
			logger.Errorf("Error scanning role for user %s: %v", id, err)
			return nil, err
		}
		roles = append(roles, &role)
	}

	var roleDtos []*rbac.RbacRoleDto
	var permissions []rbac.RbacPermission
	for _, role := range roles {
		rows, err := s.store.Query(ctx, JOIN_RBAC_ROLE_PERMISSIONS, role.ID)
		if err != nil {
			logger.Errorf("Error fetching permissions for role %s: %v", role.Name, err)
			if pgxscan.NotFound(err) {
				return nil, storage.ErrKeyNotFound
			}
			return nil, err
		}
		defer rows.Close()
		rs := pgxscan.NewRowScanner(rows)

		for rows.Next() {
			var perm rbac.RbacPermission
			if err := rs.Scan(&perm); err != nil {
				logger.Errorf("Error scanning permission for role %s: %v", role.Name, err)
				return nil, err
			}
			permissions = append(permissions, perm)
		}

		roleDtos = append(roleDtos, &rbac.RbacRoleDto{
			RbacRole:    *role,
			Permissions: permissions,
		})
	}

	return roleDtos, nil
}

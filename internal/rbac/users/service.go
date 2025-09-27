package users

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"

	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/rbac"
	"github.com/julianstephens/feature-flag-service/internal/storage"
	"github.com/julianstephens/go-utils/logger"
)

// TODO: standardize role schema
type RbacUserService struct {
	conf  *config.Config
	store *storage.PostgresStore
}

const (
	RBAC_USER_TABLE       = "rbac_users"
	RBAC_USER_ROLES_TABLE = "rbac_user_roles"

	JOIN_RBAC_USER_ROLES = `SELECT r.name
		FROM rbac_roles r
		JOIN rbac_user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1;`
)

func NewRbacUserService(conf *config.Config, store *storage.PostgresStore) *RbacUserService {
	return &RbacUserService{
		conf:  conf,
		store: store,
	}
}

func (s *RbacUserService) CreateUser(username, email, password string, roles []string) (string, error) {
	// id := uuid.New().String()

	// if err := s.store.Post(context.Background(), id, map[string]interface{}{
	// 	"username": username,
	// 	"email":    email,
	// 	"password": password,
	// 	"roles":    roles,
	// }, s.rbacUserTable); err != nil {
	// 	return "", err
	// }

	return "", nil
}

func (s *RbacUserService) GetUserByEmail(email string) (*rbac.RbacUserDto, error) {
	var rbacUser rbac.RbacUser
	if err := s.store.Get(context.Background(), RBAC_USER_TABLE, &rbacUser, "email=$1", email); err != nil {
		logger.Errorf("Error fetching user by email %s: %v", email, err)
		if pgxscan.NotFound(err) {
			return nil, storage.ErrKeyNotFound
		}
		return nil, err
	}

	var roles []string
	rows, err := s.store.Query(context.Background(), JOIN_RBAC_USER_ROLES, rbacUser.ID)
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
		var role string
		if err := rs.Scan(&role); err != nil {
			logger.Errorf("Error scanning role for user %s: %v", rbacUser.ID, err)
			return nil, err
		}
		roles = append(roles, role)
	}

	return &rbac.RbacUserDto{
		RbacUser: rbacUser,
		Roles:    roles,
	}, nil
}

func (s *RbacUserService) UpdateUser(id, username, email, password string, roles []string) error {
	// return s.store.Put(context.Background(), id, map[string]interface{}{
	// 	"username": username,
	// 	"email":    email,
	// 	"password": password,
	// 	"roles":    roles,
	// }, s.rbacUserTable)
	return nil
}

func (s *RbacUserService) DeleteUser(id string) error {
	// return s.store.Delete(context.Background(), id, s.rbacUserTable)
	return nil
}

func (s *RbacUserService) ListUsers() ([]map[string]interface{}, error) {
	// users, err := s.store.List(context.Background(), "", s.rbacUserTable)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

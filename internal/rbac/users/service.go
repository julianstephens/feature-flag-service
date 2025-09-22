package users

import (
	"context"
	"encoding/json"

	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/rbac"
	"github.com/julianstephens/feature-flag-service/internal/storage"
)

// TODO: standardize role schema
type RbacUserService struct {
	conf  *config.Config
	store storage.Store[any]
}

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

func (s *RbacUserService) GetUserByEmail(email string) (*rbac.RbacUser, error) {
	user, err := s.store.Get(context.Background(), email)
	if err != nil {
		return nil, err
	}

	var rbacUser rbac.RbacUser
	if err := json.Unmarshal([]byte(user), &rbacUser); err != nil {
		return nil, err
	}

	return &rbacUser, nil
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

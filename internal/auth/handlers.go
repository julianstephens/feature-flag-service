package auth

import (
	"errors"
	"net/http"

	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
)

var (
	ErrUserNotFound = errors.New("No user associated with email address")
)

func LoginHandler(authSvc *AuthClient, rbacUserService *users.RbacUserService) http.HandlerFunc {
	return authutils.AuthenticationHandler(authSvc.Manager, func(username, password string) (*authutils.UserInfo, error) {
		rbacUser, err := rbacUserService.GetUserByEmail(username)
		if err != nil {
			return nil, err
		}

		if !authutils.CheckPasswordHash(password, rbacUser.Password) {
			return nil, errors.New("invalid password")
		}

		return &authutils.UserInfo{
			UserID:   rbacUser.ID,
			Username: rbacUser.Email,
			Email:    rbacUser.Email,
			Roles:    rbacUser.Roles,
		}, nil
	})
}

func RefreshHandler(authSvc *AuthClient) http.HandlerFunc {
	return authutils.RefreshTokenHandler(authSvc.Manager)
}

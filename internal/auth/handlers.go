package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/security"
)

var (
	ErrUserNotFound = errors.New("no user associated with email address")
)

func LoginHandler(authSvc *AuthClient, rbacUserService *users.RbacUserService) http.HandlerFunc {
	return authutils.AuthenticationHandler(authSvc.Manager, func(username, password string) (*authutils.UserInfo, error) {
		ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		rbacUser, err := rbacUserService.GetUserByEmail(ctx, username)
		if err != nil {
			return nil, err
		}

		if !security.VerifyPassword(password, rbacUser.Password) {
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

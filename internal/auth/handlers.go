package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/feature-flag-service/internal/storage"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/httputil/request"
	"github.com/julianstephens/go-utils/httputil/response"
	"github.com/julianstephens/go-utils/security"
	"github.com/julianstephens/go-utils/validator"
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

func ActivateHandler(authSvc *AuthClient) http.HandlerFunc {
	responder := response.NewEmpty()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		var req ActivateRequest
		if err := request.DecodeJSON(r, &req); err != nil {
			responder.BadRequest(w, r, err, nil)
			return
		}

		if err := validator.ValidateEmail(req.Email); err != nil {
			responder.BadRequest(w, r, errors.New("invalid email address"), &map[string]any{"error": err.Error()})
			return
		}

		if err := validator.ValidateNonEmpty(req.Password); err != nil {
			responder.BadRequest(w, r, errors.New("password is required"), &map[string]any{"error": err.Error()})
			return
		}

		if err := validator.ValidatePassword(req.NewPassword); err != nil {
			responder.BadRequest(w, r, errors.New("invalid password"), &map[string]any{"error": err.Error()})
			return
		}

		resp, err := authSvc.Activate(ctx, req.Email, req.Password, req.NewPassword)
		if err != nil {
			if errors.Is(err, storage.ErrKeyNotFound) {
				responder.NotFound(w, r, errors.New("user not found"), nil)
				return
			}
			responder.InternalServerError(w, r, err, nil)
			return
		}

		responder.OK(w, r, resp)
	}
}

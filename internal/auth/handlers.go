package auth

import (
	"errors"
	"net/http"

	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/feature-flag-service/internal/storage"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/httputil/request"
	"github.com/julianstephens/go-utils/httputil/response"
	"github.com/julianstephens/go-utils/logger"
)

var (
	ErrUserNotFound = errors.New("No user associated with email address")
)

func LoginHandler(authSvc *AuthClient, rbacUserService *users.RbacUserService, responder *response.Responder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := request.DecodeJSON(r, &req); err != nil {
			logger.Errorf("Failed to decode login request: %v", err)
			responder.BadRequest(w, r, err)
			return
		}

		rbacUser, err := rbacUserService.GetUserByEmail(req.Email)
		if err != nil {
			if errors.Is(err, storage.ErrKeyNotFound) {
				responder.Unauthorized(w, r, ErrUserNotFound)
				return
			}
			responder.Unauthorized(w, r, err)
			return
		}

		if !authutils.CheckPasswordHash(req.Password, rbacUser.Password) {
			responder.Unauthorized(w, r, "Invalid password")
			return
		}

		token, err := authSvc.Issue(req.Email, rbacUser.Roles, &map[string]any{
			"name": rbacUser.Name,
			"uid":  rbacUser.ID,
		})
		if err != nil || token == nil {
			responder.InternalServerError(w, r, "Failed to issue token")
			return
		}

		responder.OK(w, r, token)
	}
}

func RefreshHandler(authSvc *AuthClient, rbacUserService *users.RbacUserService, responder *response.Responder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest
		if err := request.DecodeJSON(r, &req); err != nil {
			logger.Errorf("Failed to decode refresh request: %v", err)
			responder.BadRequest(w, r, "Unable to parse request body")
			return
		}

		token, err := authSvc.Refresh(req.RefreshToken)
		if err != nil {
			responder.Unauthorized(w, r, "Unable to refresh token")
			return
		}

		responder.OK(w, r, token)
	}
}

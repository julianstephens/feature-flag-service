package routes

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/rbac"
	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/httputil/response"
)

func RegisterRbacUserRoutes(router *mux.Router, userSvc *users.RbacUserService, authSvc *auth.AuthClient, responder *response.Responder) {
	router.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		user, err := getUserFromRequest(ctx, r, userSvc, authSvc)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}
		responder.OK(w, r, user)
	}).Methods("GET")
	router.HandleFunc("/me/roles", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		user, err := getUserFromRequest(ctx, r, userSvc, authSvc)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}

		res, err := userSvc.ListUserRoles(ctx, user.ID)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}
		responder.OK(w, r, res)
	}).Methods("GET")
}

func getUserFromRequest(ctx context.Context, r *http.Request, userSvc *users.RbacUserService, authSvc *auth.AuthClient) (user *rbac.RbacUserDto, err error) {
	var token string
	token, err = authutils.ExtractTokenFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		return
	}

	var claims *authutils.Claims
	claims, err = authSvc.Manager.ValidateToken(token)
	if err != nil {
		return
	}

	id := claims.UserID

	user, err = userSvc.GetUser(ctx, id)
	if err != nil {
		return
	}

	return
}

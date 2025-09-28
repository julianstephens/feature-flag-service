package routes

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	"github.com/julianstephens/go-utils/httputil/response"
)

func RegisterRbacUserRoutes(router *mux.Router, userSvc *users.RbacUserService, authSvc *auth.AuthClient, responder *response.Responder) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		res, err := userSvc.ListUsers(ctx)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}
		responder.OK(w, r, res)
	}).Methods("GET")
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		id := mux.Vars(r)["id"]

		res, err := userSvc.GetUser(ctx, id)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}
		responder.OK(w, r, res)
	}).Methods("GET")
	router.HandleFunc("/{id}/roles", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		id := mux.Vars(r)["id"]

		res, err := userSvc.ListUserRoles(ctx, id)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}
		responder.OK(w, r, res)
	}).Methods("GET")
}

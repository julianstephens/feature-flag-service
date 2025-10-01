package routes

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/flag"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	"github.com/julianstephens/go-utils/httputil/middleware"
	"github.com/julianstephens/go-utils/httputil/request"
	"github.com/julianstephens/go-utils/httputil/response"
)

func RegisterFlagRoutes(router *mux.Router, flagSvc flag.Service, authSvc *auth.AuthClient, responder *response.Responder) {
	readRoutes := router.PathPrefix("").Subrouter()
	readRoutes.Use(middleware.RequireRoles(authSvc.Manager, "user", "editor", "admin"))
	readRoutes.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		res, err := flagSvc.ListFlags(ctx)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}
		responder.OK(w, r, res)
	}).Methods("GET")
	readRoutes.HandleFunc("/{flagKey}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		vars := mux.Vars(r)
		flagKey := vars["flagKey"]

		res, err := flagSvc.GetFlag(ctx, flagKey)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}

		responder.OK(w, r, res)
	}).Methods("GET")

	writeRoutes := router.PathPrefix("").Subrouter()
	writeRoutes.Use(middleware.RequireRoles(authSvc.Manager, "editor", "admin"))
	writeRoutes.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
		defer cancel()

		var req ffpb.CreateFlagRequest
		if err := request.DecodeJSON(r, &req); err != nil {
			responder.BadRequest(w, r, err, nil)
			return
		}

		res, err := flagSvc.CreateFlag(ctx, req.Name, req.Description, req.Enabled)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}
		responder.Created(w, r, res)
	}).Methods("POST")
	writeRoutes.HandleFunc("/{flagKey}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
		defer cancel()
		vars := mux.Vars(r)
		flagKey := vars["flagKey"]

		var req ffpb.UpdateFlagRequest
		if err := request.DecodeJSON(r, &req); err != nil {
			responder.BadRequest(w, r, err, nil)
			return
		}

		res, err := flagSvc.UpdateFlag(ctx, flagKey, req.Name, req.Description, req.Enabled)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}

		responder.OK(w, r, res)
	}).Methods("PUT")

	deleteRoutes := router.PathPrefix("").Subrouter()
	deleteRoutes.Use(middleware.RequireRoles(authSvc.Manager, "admin"))
	deleteRoutes.HandleFunc("/{flagKey}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
		defer cancel()
		vars := mux.Vars(r)
		flagKey := vars["flagKey"]

		err := flagSvc.DeleteFlag(ctx, flagKey)
		if err != nil {
			HandleError(responder, w, r, err)
			return
		}
		responder.NoContent(w, r)
	}).Methods("DELETE")
}

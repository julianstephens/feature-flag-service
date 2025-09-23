package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"google.golang.org/grpc"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/flag"
	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/go-utils/httputil/middleware"
	"github.com/julianstephens/go-utils/httputil/request"
	"github.com/julianstephens/go-utils/httputil/response"
)

const (
	DEFAULT_TIMEOUT = 30 * time.Second
)

func StartREST(addr string, conf *config.Config, services ...any) error {
	logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	responder := response.NewEmpty()

	router := mux.NewRouter()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logging(logger))
	router.Use(middleware.Recovery(errorLogger))
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	apiGrp := router.PathPrefix("/api/" + conf.APIVersion).Subrouter()
	apiGrp.HandleFunc("/checkhealth", func(w http.ResponseWriter, r *http.Request) {
		responder.OK(w, r, map[string]string{"status": "OK", "version": "1.0", "name": "Feature Flag Service"})
	})

	servicesMap := make(map[string]any)
	for _, svc := range services {
		switch s := svc.(type) {
		case flag.Service:
			servicesMap["flagService"] = s
		case *auth.AuthClient:
			servicesMap["authService"] = s
		case *users.RbacUserService:
			servicesMap["userService"] = s
		// case config.Service: --- IGNORE ---
		// 	servicesMap["configService"] = s --- IGNORE ---
		// case audit.Service: --- IGNORE ---
		// 	servicesMap["auditService"] = s --- IGNORE ---
		// case rbac.Service: --- IGNORE ---
		// 	servicesMap["rbacService"] = s --- IGNORE ---
		default:
			log.Printf("Warning: Unknown service type %T provided to StartREST", s)
		}
	}

	flagSvc := servicesMap["flagService"].(flag.Service)
	flags := apiGrp.PathPrefix("/flags").Subrouter()
	flags.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
		defer cancel()

		res, err := flagSvc.ListFlags(ctx)
		if err != nil {
			handleError(responder, w, r, err)
			return
		}
		responder.OK(w, r, res)
	}).Methods("GET")
	flags.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
		defer cancel()

		var req ffpb.CreateFlagRequest
		if err := request.DecodeJSON(r, &req); err != nil {
			responder.BadRequest(w, r, err)
			return
		}

		res, err := flagSvc.CreateFlag(ctx, req.Name, req.Description, req.Enabled)
		if err != nil {
			handleError(responder, w, r, err)
			return
		}
		responder.Created(w, r, res)
	}).Methods("POST")
	flags.HandleFunc("/{flagKey}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
		defer cancel()

		vars := mux.Vars(r)
		flagKey := vars["flagKey"]

		res, err := flagSvc.GetFlag(ctx, flagKey)
		if err != nil {
			handleError(responder, w, r, err)
			return
		}

		responder.OK(w, r, res)
	}).Methods("GET")
	flags.HandleFunc("/{flagKey}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
		defer cancel()
		vars := mux.Vars(r)
		flagKey := vars["flagKey"]

		var req ffpb.UpdateFlagRequest
		if err := request.DecodeJSON(r, &req); err != nil {
			responder.BadRequest(w, r, err)
			return
		}

		res, err := flagSvc.UpdateFlag(ctx, flagKey, req.Name, req.Description, req.Enabled)
		if err != nil {
			handleError(responder, w, r, err)
			return
		}

		responder.OK(w, r, res)
	}).Methods("PUT")
	flags.HandleFunc("/{flagKey}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
		defer cancel()
		vars := mux.Vars(r)
		flagKey := vars["flagKey"]

		err := flagSvc.DeleteFlag(ctx, flagKey)
		if err != nil {
			handleError(responder, w, r, err)
			return
		}
		responder.NoContent(w, r)
	}).Methods("DELETE")

	authSvc := servicesMap["authService"].(*auth.AuthClient)
	userSvc := servicesMap["userService"].(*users.RbacUserService)
	authGrp := apiGrp.PathPrefix("/auth").Subrouter()
	authGrp.HandleFunc("/login", auth.LoginHandler(authSvc, userSvc, responder)).Methods("POST")
	authGrp.HandleFunc("/refresh", auth.RefreshHandler(authSvc, userSvc, responder)).Methods("POST")

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return srv.ListenAndServe()
}

func RegisterGRPC(grpcServer *grpc.Server, flagSvc flag.Service) {
	ffpb.RegisterFlagServiceServer(grpcServer, &flag.FlagGRPCServer{
		UnimplementedFlagServiceServer: ffpb.UnimplementedFlagServiceServer{},
		Service:                        flagSvc,
	})
}

func handleError(responder *response.Responder, w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case context.Canceled:
		responder.ErrorWithStatus(w, r, http.StatusRequestTimeout, err)
	case context.DeadlineExceeded:
		responder.ErrorWithStatus(w, r, http.StatusRequestTimeout, err)
	case rpctypes.ErrEmptyKey:
		responder.BadRequest(w, r, err)
	default:
		responder.InternalServerError(w, r, err)
	}
}

package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/flag"
	grpcmiddleware "github.com/julianstephens/feature-flag-service/internal/grpc"
	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/feature-flag-service/internal/server/routes"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/httputil/middleware"
	"github.com/julianstephens/go-utils/httputil/response"
)

func StartREST(addr string, conf *config.Config, services ...any) error {
	logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

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
	authSvc := servicesMap["authService"].(*auth.AuthClient)
	userSvc := servicesMap["userService"].(*users.RbacUserService)
	flagSvc := servicesMap["flagService"].(flag.Service)

	responder := response.NewEmpty()

	router := mux.NewRouter()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logging(logger))
	router.Use(middleware.Recovery(errorLogger))
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	publicGrp := router.PathPrefix("/api/" + conf.APIVersion).Subrouter()
	publicGrp.HandleFunc("/checkhealth", func(w http.ResponseWriter, r *http.Request) {
		responder.OK(w, r, map[string]string{"status": "OK", "version": "1.0", "name": "Feature Flag Service"})
	})
	authGrp := publicGrp.PathPrefix("/auth").Subrouter()
	authGrp.HandleFunc("/login", auth.LoginHandler(authSvc, userSvc)).Methods("POST")
	authGrp.HandleFunc("/refresh", auth.RefreshHandler(authSvc)).Methods("POST")

	privateGroup := publicGrp.PathPrefix("").Subrouter()
	privateGroup.Use(middleware.JWTAuth(authSvc.Manager))

	routes.RegisterFlagRoutes(privateGroup.PathPrefix("/flags").Subrouter(), flagSvc, authSvc, responder)

	rbacGrp := privateGroup.PathPrefix("/rbac").Subrouter()
	routes.RegisterRbacUserRoutes(rbacGrp.PathPrefix("/users").Subrouter(), userSvc, authSvc, responder)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return srv.ListenAndServe()
}

func RegisterGRPC(grpcServer *grpc.Server, flagSvc flag.Service, authSvc auth.Service, userSvc users.Service) {
	ffpb.RegisterFlagServiceServer(grpcServer, &flag.FlagGRPCServer{
		UnimplementedFlagServiceServer: ffpb.UnimplementedFlagServiceServer{},
		Service:                        flagSvc,
	})
	ffpb.RegisterAuthServiceServer(grpcServer, &auth.AuthGRPCServer{
		UnimplementedAuthServiceServer: ffpb.UnimplementedAuthServiceServer{},
		Service:                        authSvc,
	})
	ffpb.RegisterRbacUserServiceServer(grpcServer, &users.RbacUserGRPCServer{
		UnimplementedRbacUserServiceServer: ffpb.UnimplementedRbacUserServiceServer{},
		Service:                            userSvc,
	})
}

// NewGRPCServerWithAuth creates a gRPC server with RBAC interceptors
func NewGRPCServerWithAuth(jwtManager *authutils.JWTManager) *grpc.Server {
	// Create server with method-based auth interceptor
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.MethodBasedAuthInterceptor(jwtManager, grpcmiddleware.MethodRoles)),
	}

	return grpc.NewServer(opts...)
}

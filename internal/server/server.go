package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/flag"
)

func StartREST(addr string) error {
	router := mux.NewRouter()
	router.HandleFunc("/checkhealth", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("Starting REST server on %s", addr)
	return srv.ListenAndServe()
}

type flagGRPCServer struct {
	ffpb.UnimplementedFlagServiceServer
	svc flag.Service
}

func RegisterGRPC(grpcServer *grpc.Server, flagSvc flag.Service) {
	ffpb.RegisterFlagServiceServer(grpcServer, &flagGRPCServer{
		UnimplementedFlagServiceServer: ffpb.UnimplementedFlagServiceServer{},
		svc: flagSvc,
	})
}

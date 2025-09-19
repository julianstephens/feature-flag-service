package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/flag"
	"github.com/julianstephens/go-utils/httputil/response"
)


func StartREST(addr string) error {
	responder := response.NewWithLogging()
	router := mux.NewRouter()
	router.HandleFunc("/checkhealth", func(w http.ResponseWriter, r *http.Request) {
		responder.OK(w, r, map[string]string{"status": "OK", "version": "1.0", "name": "Feature Flag Service"})
	})

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

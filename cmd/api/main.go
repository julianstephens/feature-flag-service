package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/flag"
	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	"github.com/julianstephens/feature-flag-service/internal/server"
	"github.com/julianstephens/feature-flag-service/internal/storage"
	"github.com/julianstephens/go-utils/logger"
)

func main() {
	conf := config.LoadConfig()
	etcdClient, err := storage.NewEtcdStore([]string{conf.StorageEndpoint}, "/featureflags/")
	if err != nil {
		logger.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	pgClient, err := storage.NewPostgresStore(conf)
	if err != nil {
		logger.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pgClient.Close()

	flagService := flag.NewService(conf, etcdClient)
	authService, err := auth.NewAuthClient(conf)
	if err != nil {
		logger.Fatalf("Failed to create auth service: %v", err)
	}
	userService := users.NewRbacUserService(conf, pgClient)

	go func() {
		logger.Infof("Starting REST API on :%s...", conf.HTTPPort)
		if err := server.StartREST(":"+conf.HTTPPort, conf, flagService, authService, userService); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("REST server error: %v", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:"+conf.GRPCPort)
		if err != nil {
			logger.Fatalf("Failed to listen on gRPC port %s: %v", conf.GRPCPort, err)
		}
		grpcServer := grpc.NewServer()
		server.RegisterGRPC(grpcServer, flagService)
		logger.Infof("Starting gRPC API on :%s...", conf.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("gRPC server error: %v", err)
		}
	}()

	waitForShutdown()
	logger.Info("API service stopped.")
}

func waitForShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logger.Info("Shutdown signal received, exiting...")
}

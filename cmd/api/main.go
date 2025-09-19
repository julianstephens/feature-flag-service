package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/flag"
	"github.com/julianstephens/feature-flag-service/internal/server"
	"github.com/julianstephens/feature-flag-service/internal/storage"
)

func main() {
	conf := config.LoadConfig()
	etcdClient, err := storage.NewEtcdStore([]string{conf.StorageEndpoint}, "/featureflags/")
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Client.Close()
	flagService := flag.NewService(conf, etcdClient)

	go func() {
		log.Printf("Starting REST API on :%s...", conf.HTTPPort)
		if err := server.StartREST(":" + conf.HTTPPort, conf, flagService); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST server error: %v", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:"+conf.GRPCPort)
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port %s: %v", conf.GRPCPort, err)
		}
		grpcServer := grpc.NewServer()
		server.RegisterGRPC(grpcServer, flagService)
		log.Printf("Starting gRPC API on :%s...", conf.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	waitForShutdown()
	log.Println("API service stopped.")
}

func waitForShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutdown signal received, exiting...")
}

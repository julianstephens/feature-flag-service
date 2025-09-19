package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/server"
)

func main() {
	conf := config.LoadConfig()
	// flagService := flag.NewService(conf)

	go func() {
		log.Printf("Starting REST API on :%s...", conf.HTTPPort)
		if err := server.StartREST(":" + conf.HTTPPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST server error: %v", err)
		}
	}()

	// go func() {
	// 	lis, err := net.Listen("tcp", ":"+grpcPort)
	// 	if err != nil {
	// 		log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
	// 	}
	// 	grpcServer := grpc.NewServer()
	// 	server.RegisterGRPC(grpcServer, flagService, configService, auditService, rbacService)
	// 	log.Printf("Starting gRPC API on :%s...", grpcPort)
	// 	if err := grpcServer.Serve(lis); err != nil {
	// 		log.Fatalf("gRPC server error: %v", err)
	// 	}
	// }()
	waitForShutdown()
	log.Println("API service stopped.")
}

func waitForShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutdown signal received, exiting...")
}

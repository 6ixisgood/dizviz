package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/6ixisgood/matrix-ticker/pkg/api"
	_ "github.com/6ixisgood/matrix-ticker/pkg/component"
	"github.com/6ixisgood/matrix-ticker/pkg/controlplane"
	_ "github.com/6ixisgood/matrix-ticker/pkg/view/types"
)

func main() {
	// Parse flags
	configFilePath := flag.String("config", "./config/controlplane.yaml", "Path to yaml config file")
	flag.Parse()

	// Load configuration
	LoadConfig(configFilePath)

	// Use config values or defaults from flags
	grpcAddr := "localhost:" + Config.Server.GRPCPort
	httpAddr := "localhost:" + Config.Server.HTTPPort
	if Config.Server.GRPCPort == "" {
		grpcAddr = "localhost:50051"
	}
	if Config.Server.HTTPPort == "" {
		httpAddr = "localhost:8080"
	}

	log.Println("=== DizViz Control Plane ===")
	log.Printf("gRPC Server: %s", grpcAddr)
	log.Printf("HTTP API: %s", httpAddr)

	log.Println("[ControlPlane] System initialized")

	// Create gRPC server
	server := controlplane.NewServer(grpcAddr)

	// Set up HTTP API with registry and server
	api.SetRegistry(server.GetRegistry())
	api.SetServer(server)
	router := api.Router()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nShutdown signal received")
		server.Stop()
		os.Exit(0)
	}()

	// Start HTTP API server in background
	go func() {
		log.Printf("[HTTP API] Starting on %s", httpAddr)
		if err := router.Run(httpAddr); err != nil {
			log.Printf("[HTTP API] Server error: %v", err)
		}
	}()

	// Start gRPC server (blocks until stopped)
	if err := server.Start(); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

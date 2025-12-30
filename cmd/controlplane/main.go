package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/6ixisgood/matrix-ticker/pkg/api"
	"github.com/6ixisgood/matrix-ticker/pkg/controlplane"
)

func main() {
	// Parse flags
	grpcAddr := flag.String("grpc", "localhost:50051", "gRPC server address")
	httpAddr := flag.String("http", "localhost:8080", "HTTP API server address")
	flag.Parse()

	log.Println("=== DizViz Control Plane ===")
	log.Printf("gRPC Server: %s", *grpcAddr)
	log.Printf("HTTP API: %s", *httpAddr)

	// Create gRPC server
	server := controlplane.NewServer(*grpcAddr)

	// Set up HTTP API with registry
	api.SetRegistry(server.GetRegistry())
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
		log.Printf("[HTTP API] Starting on %s", *httpAddr)
		if err := router.Run(*httpAddr); err != nil {
			log.Printf("[HTTP API] Server error: %v", err)
		}
	}()

	// Start gRPC server (blocks until stopped)
	if err := server.Start(); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

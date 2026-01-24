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
	"github.com/6ixisgood/matrix-ticker/pkg/store"
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

	// Initialize store
	storeDir := Config.Data.StoreDir
	if storeDir == "" {
		storeDir = "./data/store" // Default fallback
	}
	st, err := store.NewStore(storeDir)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}
	log.Printf("[ControlPlane] Store initialized at: %s", storeDir)

	// Create store service
	storeService := controlplane.NewStoreService(st)
	defer storeService.Close()

	log.Println("[ControlPlane] System initialized")

	// Build data source config from loaded config
	dataSourceCfg := &controlplane.DataSourceConfig{
		Sleeper: struct {
			BaseUrl string
		}{
			BaseUrl: Config.Data.Sleeper.BaseUrl,
		},
		SportsFeed: struct {
			BaseUrl  string
			Username string
			Password string
		}{
			BaseUrl:  Config.Data.SportsFeed.BaseUrl,
			Username: Config.Data.SportsFeed.Username,
			Password: Config.Data.SportsFeed.Password,
		},
		Weather: struct {
			BaseUrl string
			Key     string
		}{
			BaseUrl: Config.Data.Weather.BaseUrl,
			Key:     Config.Data.Weather.Key,
		},
	}

	// Create gRPC server with store service and data source config
	server := controlplane.NewServer(grpcAddr, storeService, dataSourceCfg)

	// Set up HTTP API with registry, server, and store service
	api.SetRegistry(server.GetRegistry())
	api.SetServer(server)
	api.SetStoreService(storeService)
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

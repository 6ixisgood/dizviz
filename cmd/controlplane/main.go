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
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	_ "github.com/6ixisgood/matrix-ticker/pkg/view/types"
)

func main() {
	// Parse flags
	configFilePath := flag.String("config", "./config/config.yaml", "Path to yaml config file")
	grpcAddr := flag.String("grpc", "localhost:50051", "gRPC server address")
	httpAddr := flag.String("http", "localhost:8080", "HTTP API server address")
	flag.Parse()

	// Load configuration
	LoadConfig(configFilePath)

	log.Println("=== DizViz Control Plane ===")
	log.Printf("gRPC Server: %s", *grpcAddr)
	log.Printf("HTTP API: %s", *httpAddr)
	log.Printf("Store Directory: %s", AppConfig.Data.StoreDir)

	// Initialize the store for view definitions
	// appStore, err := store.NewStore(AppConfig.Data.StoreDir)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize store: %v", err)
	// }
	// defer appStore.Close()

	// Configure viewCommon so view definitions can be loaded
	viewCommon.SetViewCommonConfig(&viewCommon.ViewCommonConfig{
		MatrixRows:        AppConfig.Matrix.Rows * AppConfig.Matrix.Parallel,
		MatrixCols:        AppConfig.Matrix.Cols * AppConfig.Matrix.Chain,
		ImageDir:          AppConfig.Data.ImageDir,
		CacheDir:          AppConfig.Data.CacheDir,
		DefaultImageSizeX: AppConfig.Default.ImageSizeX,
		DefaultImageSizeY: AppConfig.Default.ImageSizeY,
		DefaultFontSize:   AppConfig.Default.FontSize,
		DefaultFontColor:  AppConfig.Default.FontColor,
		DefaultFontStyle:  AppConfig.Default.FontStyle,
		DefaultFontType:   AppConfig.Default.FontType,
		// Store:             appStore,
	})

	// Configure utils
	util.SetUtilConfig(&util.UtilConfig{
		CacheDir: AppConfig.Data.CacheDir,
		FontDir:  AppConfig.Data.FontDir,
	})

	log.Println("[ControlPlane] View system initialized")

	// Create gRPC server
	server := controlplane.NewServer(*grpcAddr)

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

package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/6ixisgood/matrix-ticker/pkg/agent"
	_ "github.com/6ixisgood/matrix-ticker/pkg/component/types"
	"github.com/6ixisgood/matrix-ticker/pkg/display"
	"github.com/6ixisgood/matrix-ticker/pkg/store"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	_ "github.com/6ixisgood/matrix-ticker/pkg/view/types"
	rgbmatrix "github.com/sixisgoood/go-rpi-rgb-led-matrix"
)

var (
	configFilePath = flag.String("config", "./config/agent.yaml", "path to yaml config file")
)

func main() {
	flag.Parse()

	// Load agent configuration
	LoadConfig(configFilePath)

	log.Printf("=== DizViz Agent ===")
	log.Printf("Agent Name: %s", Config.Agent.Name)
	log.Printf("Control Plane: %s", Config.Agent.ControlPlaneAddr)
	log.Printf("Displays: %d configured", len(Config.Displays))

	// Ensure at least one display is configured
	if len(Config.Displays) == 0 {
		log.Fatalf("No displays configured in config file")
	}

	// For now, we only support the first display
	// TODO: Support multiple displays
	displayConfig := Config.Displays[0]

	log.Printf("Using display: %s (type: %s)", displayConfig.ID, displayConfig.Type)

	// Create display based on type
	var disp display.Display
	var matrixConfig *rgbmatrix.HardwareConfig

	switch displayConfig.Type {
	case "matrix":
		matrixConfig = &rgbmatrix.DefaultConfig
		matrixConfig.Rows = displayConfig.Hardware.Rows
		matrixConfig.Cols = displayConfig.Hardware.Cols
		matrixConfig.Parallel = displayConfig.Hardware.Parallel
		matrixConfig.ChainLength = displayConfig.Hardware.Chain
		matrixConfig.Brightness = displayConfig.Hardware.Brightness
		matrixConfig.HardwareMapping = displayConfig.Hardware.HardwareMapping
		matrixConfig.ShowRefreshRate = displayConfig.Hardware.ShowRefresh
		matrixConfig.InverseColors = displayConfig.Hardware.InverseColors
		matrixConfig.DisableHardwarePulsing = displayConfig.Hardware.DisableHardwarePulsing
		matrixConfig.GpioSlowdown = displayConfig.Hardware.GpioSlowdown
		matrixConfig.RateLimitHz = displayConfig.Hardware.RateLimitHz

		disp = display.NewMatrixDisplay()
		log.Printf("LED Matrix display initialized: %dx%d",
			displayConfig.Hardware.Cols*displayConfig.Hardware.Chain,
			displayConfig.Hardware.Rows*displayConfig.Hardware.Parallel)

	case "file":
		log.Fatalf("File display not yet implemented")

	case "hdmi":
		log.Fatalf("HDMI display not yet implemented")

	default:
		log.Fatalf("Unknown display type: %s", displayConfig.Type)
	}

	// Create agent configuration
	agentConfig := agent.Config{
		AgentName:           Config.Agent.Name,
		ControlPlaneAddr:    Config.Agent.ControlPlaneAddr,
		HeartbeatInterval:   Config.Agent.HeartbeatInterval,
		ReconnectDelay:      Config.Agent.ReconnectDelay,
		EnableAutoReconnect: Config.Agent.AutoReconnect,
		Version:             Config.Agent.Version,
		DisplayType:         displayConfig.Type,
		FPS:                 displayConfig.FPS,
		BufferSize:          displayConfig.BufferSize,
	}

	// Create capabilities from display
	capabilities := agent.NewCapabilitiesFromDisplay(
		disp,
		displayConfig.ID,
		displayConfig.Type,
		agent.DefaultSupportedViews(),
	)

	// Create agent
	ag := agent.New(agentConfig, capabilities)

	// Initialize store
	st, err := store.NewStore(Config.Runtime.CacheDir)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Configure utils for font/cache handling
	util.SetUtilConfig(&util.UtilConfig{
		CacheDir: Config.Runtime.CacheDir,
		FontDir:  Config.Runtime.FontsDir,
	})

	// Convert data sources to the expected format
	dataSources := make(map[string]interface{})
	for key, value := range Config.Runtime.DataSources {
		dataSources[key] = value
	}

	// Build agent context
	agentCtx := &viewCommon.AgentContext{
		ImageDir:    Config.Runtime.ImagesDir,
		CacheDir:    Config.Runtime.CacheDir,
		FontsDir:    Config.Runtime.FontsDir,
		Store:       st,
		DataSources: dataSources,
	}
	ag.SetAgentContext(agentCtx)

	// Calculate display dimensions
	matrixRows := displayConfig.Hardware.Rows * displayConfig.Hardware.Parallel
	matrixCols := displayConfig.Hardware.Cols * displayConfig.Hardware.Chain

	// Prepare font defaults from config
	fontDefaults := viewCommon.DisplayContext{
		DefaultFontSize:  Config.Runtime.Defaults.FontSize,
		DefaultFontColor: Config.Runtime.Defaults.FontColor,
		DefaultFontStyle: Config.Runtime.Defaults.FontStyle,
		DefaultFontType:  Config.Runtime.Defaults.FontType,
	}

	// Add the display to the agent
	if err := ag.AddDisplay(displayConfig.ID, disp, displayConfig.FPS, displayConfig.BufferSize, matrixRows, matrixCols, fontDefaults); err != nil {
		log.Fatalf("Failed to add display to agent: %v", err)
	}

	// Set initial view (welcome message) BEFORE starting the agent
	welcomeConfig := []byte(`{
		"text": "DizViz Agent Ready",
		"alignment": "center",
		"justify": "center",
		"color": "#00FF00FF",
		"bg-color": "#002288FF"
	}`)

	regView := viewCommon.RegisteredViews["text"]
	configInstance := regView.NewConfig()
	if err := json.Unmarshal(welcomeConfig, &configInstance); err != nil {
		log.Fatalf("Failed to unmarshal welcome view config: %v", err)
	}

	welcomeView, err := regView.NewView(configInstance)
	if err != nil {
		log.Fatalf("Failed to create welcome view: %v", err)
	}

	// Inject context into welcome view using defaults from config
	displayCtx := &viewCommon.DisplayContext{
		MatrixRows:        matrixRows,
		MatrixCols:        matrixCols,
		DefaultImageSizeX: matrixCols,
		DefaultImageSizeY: matrixRows,
		DefaultFontSize:   Config.Runtime.Defaults.FontSize,
		DefaultFontColor:  Config.Runtime.Defaults.FontColor,
		DefaultFontStyle:  Config.Runtime.Defaults.FontStyle,
		DefaultFontType:   Config.Runtime.Defaults.FontType,
	}
	welcomeViewCtx := &viewCommon.ViewContext{
		Agent:   agentCtx,
		Display: displayCtx,
	}
	welcomeView.SetContext(welcomeViewCtx)

	if err := ag.SetInitialView(displayConfig.ID, welcomeView); err != nil {
		log.Fatalf("Failed to set initial view: %v", err)
	}

	// Now start agent (will initialize the view properly)
	if err := ag.Start(matrixConfig); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	log.Printf("Agent started successfully!")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Printf("Shutdown signal received")
	ag.Stop()
	log.Printf("Agent stopped gracefully")
}

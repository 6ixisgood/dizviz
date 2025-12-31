package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v2"

	"github.com/6ixisgood/matrix-ticker/pkg/agent"
	_ "github.com/6ixisgood/matrix-ticker/pkg/component/types"
	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	"github.com/6ixisgood/matrix-ticker/pkg/display"
	"github.com/6ixisgood/matrix-ticker/pkg/store"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	_ "github.com/6ixisgood/matrix-ticker/pkg/view/types"
	rgbmatrix "github.com/sixisgoood/go-rpi-rgb-led-matrix"
)

var (
	configFilePath   = flag.String("config", "./config/config.yaml", "path to yaml config file")
	agentName        = flag.String("name", "dizviz-agent", "agent name")
	controlPlaneAddr = flag.String("control-plane", "", "control plane address (e.g., localhost:50051)")
	displayType      = flag.String("display", "led-matrix", "display type: led-matrix, hdmi, file")
)

type ApplicationConfig struct {
	Matrix struct {
		Rows                   int    `yaml:"rows"`
		Cols                   int    `yaml:"cols"`
		Parallel               int    `yaml:"parallel"`
		Chain                  int    `yaml:"chain"`
		Brightness             int    `yaml:"brightness"`
		HardwareMapping        string `yaml:"harware_mapping"`
		ShowRefresh            bool   `yaml:"show_refresh"`
		InverseColors          bool   `yaml:"inverse_colors"`
		DisableHardwarePulsing bool   `yaml:"disable_hardware_pulsing"`
		GpioSlowdown           int    `yaml:"gpio_slowdown"`
		RateLimitHz            int    `yaml:"rate_limit_hz"`
	} `yaml:"matrix"`
	Default struct {
		ImageSizeX int    `yaml:"image_size_x"`
		ImageSizeY int    `yaml:"image_size_y"`
		FontSize   int    `yaml:"font_size"`
		FontColor  string `yaml:"font_color"`
		FontStyle  string `yaml:"font_style"`
		FontType   string `yaml:"font_type"`
	}
	Server struct {
		AllowedHosts string `yaml:"allowed_hosts"`
		Port         string `yaml:"port"`
	}
	Data struct {
		ImageDir string `yaml:"images"`
		CacheDir string `yaml:"cache"`
		FontDir  string `yaml:"fonts"`
		StoreDir string `yaml:"store"`
		Sleeper  struct {
			BaseUrl string `yaml:"baseUrl"`
		}
		SportsFeed struct {
			BaseUrl  string `yaml:"baseUrl"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"sportsfeed"`
		Weather struct {
			BaseUrl string `yaml:"baseUrl"`
			Key     string `yaml:"key"`
		} `yaml:"weather"`
	} `yaml:"data"`
}

var AppConfig = ApplicationConfig{}

func loadConfig(filepath *string) {
	data, err := ioutil.ReadFile(*filepath)
	if err != nil {
		log.Fatalf("Error loading config: '%v'", err)
		return
	}

	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		log.Fatalf("Error unmarshaling app config: '%v'", err)
		return
	}
}

func main() {
	flag.Parse()

	log.Printf("=== DizViz Agent ===")
	log.Printf("Agent Name: %s", *agentName)
	log.Printf("Display Type: %s", *displayType)
	if *controlPlaneAddr != "" {
		log.Printf("Control Plane: %s", *controlPlaneAddr)
	} else {
		log.Printf("Control Plane: Standalone mode (no control plane)")
	}

	// Load application config
	loadConfig(configFilePath)

	// Initialize store
	appStore, err := store.NewStore(AppConfig.Data.StoreDir)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}
	defer appStore.Close()

	// Configure views
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
		Store:             appStore,
	})

	// Configure utils
	util.SetUtilConfig(&util.UtilConfig{
		CacheDir: AppConfig.Data.CacheDir,
		FontDir:  AppConfig.Data.FontDir,
	})

	// Initialize data sources
	d.InitSportsFeedClient(d.SportsFeedConfig{
		BaseUrl:  AppConfig.Data.SportsFeed.BaseUrl,
		Username: AppConfig.Data.SportsFeed.Username,
		Password: AppConfig.Data.SportsFeed.Password,
	})

	d.InitSleeperClient(d.SleeperConfig{
		BaseUrl: AppConfig.Data.Sleeper.BaseUrl,
	})

	// Create display based on type
	var disp display.Display
	var matrixConfig *rgbmatrix.HardwareConfig

	switch *displayType {
	case "led-matrix":
		matrixConfig = &rgbmatrix.DefaultConfig
		matrixConfig.Rows = AppConfig.Matrix.Rows
		matrixConfig.Cols = AppConfig.Matrix.Cols
		matrixConfig.Parallel = AppConfig.Matrix.Parallel
		matrixConfig.ChainLength = AppConfig.Matrix.Chain
		matrixConfig.Brightness = AppConfig.Matrix.Brightness
		matrixConfig.HardwareMapping = AppConfig.Matrix.HardwareMapping
		matrixConfig.ShowRefreshRate = AppConfig.Matrix.ShowRefresh
		matrixConfig.InverseColors = AppConfig.Matrix.InverseColors
		matrixConfig.DisableHardwarePulsing = AppConfig.Matrix.DisableHardwarePulsing
		matrixConfig.GpioSlowdown = AppConfig.Matrix.GpioSlowdown
		matrixConfig.RateLimitHz = AppConfig.Matrix.RateLimitHz

		disp = display.NewMatrixDisplay()
		log.Printf("LED Matrix display initialized: %dx%d",
			AppConfig.Matrix.Cols*AppConfig.Matrix.Chain,
			AppConfig.Matrix.Rows*AppConfig.Matrix.Parallel)

	case "file":
		// TODO: Implement file-based display for testing
		log.Fatalf("File display not yet implemented")

	case "hdmi":
		// TODO: Implement framebuffer display for HDMI
		log.Fatalf("HDMI display not yet implemented")

	default:
		log.Fatalf("Unknown display type: %s", *displayType)
	}

	// Create agent configuration
	agentConfig := agent.DefaultConfig()
	agentConfig.AgentName = *agentName
	agentConfig.ControlPlaneAddr = *controlPlaneAddr
	agentConfig.DisplayType = *displayType
	agentConfig.FPS = 30
	agentConfig.BufferSize = 20

	// Create capabilities from display
	capabilities := agent.NewCapabilitiesFromDisplay(
		disp,
		"primary",    // display ID
		*displayType, // display type
		agent.DefaultSupportedViews(),
	)

	// Create agent
	ag := agent.New(agentConfig, capabilities)

	// Add the display to the agent
	if err := ag.AddDisplay("primary", disp, agentConfig.FPS, agentConfig.BufferSize); err != nil {
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

	var welcomeViewConfig viewCommon.ViewConfig
	if err := json.Unmarshal(welcomeConfig, &welcomeViewConfig); err != nil {
		log.Fatalf("Failed to parse welcome view config: %v", err)
	}

	regView := viewCommon.RegisteredViews["text"]

	configInstance := regView.NewConfig()
	if err := json.Unmarshal(welcomeConfig, &configInstance); err != nil {
		log.Fatalf("Failed to unmarshal welcome view config: %v", err)
	}

	welcomeView, err := regView.NewView(configInstance)
	if err != nil {
		log.Fatalf("Failed to create welcome view: %v", err)
	}

	if err := ag.SetInitialView("primary", welcomeView); err != nil {
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

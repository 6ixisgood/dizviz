package main

import (
	"flag"
	"fmt"
	"log"

	"encoding/json"

	"github.com/6ixisgood/matrix-ticker/pkg/api"
	"github.com/6ixisgood/matrix-ticker/pkg/app"
	_ "github.com/6ixisgood/matrix-ticker/pkg/component"
	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	"github.com/6ixisgood/matrix-ticker/pkg/display"
	"github.com/6ixisgood/matrix-ticker/pkg/store"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	_ "github.com/6ixisgood/matrix-ticker/pkg/view/types"
)

/*

1) Read command line args
2) Load AppConfig config
3) Set the rgbmatrix config
4) Create/configure data sources
5) Init the View Config
	a) Matrix Size, Image dir, Cache dir, defaults, etc.
	b) Register all views
5) Set component config
	a) register all components
6) Create/Configure webserver
7) Start animation
8) Start webserver


*/

var (
	configFilePath = flag.String("config", "./config.yaml", "path to yaml config file")
	application    *app.Application
)

func main() {
	flag.Parse()

	LoadConfig(configFilePath)

	// set the RBG matrix configs
	matrixConfig := &DefaultConfig{}
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

	// init the store
	appStore, err := store.NewStore(AppConfig.Data.StoreDir)
	if err != nil {
		panic(err)
	}
	defer appStore.Close()

	// configure the views
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

	// configure utils
	util.SetUtilConfig(&util.UtilConfig{
		CacheDir: AppConfig.Data.CacheDir,
		FontDir:  AppConfig.Data.FontDir,
	})

	// init the sports feed client
	d.InitSportsFeedClient(d.SportsFeedConfig{
		BaseUrl:  AppConfig.Data.SportsFeed.BaseUrl,
		Username: AppConfig.Data.SportsFeed.Username,
		Password: AppConfig.Data.SportsFeed.Password,
	})

	// init the sleeper client
	d.InitSleeperClient(d.SleeperConfig{
		BaseUrl: AppConfig.Data.Sleeper.BaseUrl,
	})

	// Initialize global compositor
	// view.InitGlobalCompositor()

	// Create application
	application = app.New()

	// Create matrix display
	matrixDisplay := display.NewMatrixDisplay()

	// setup matrix
	fmt.Println("Starting Matrix\n")

	log.Printf("Initializing the starting view")

	t := "text"
	configJSON := []byte(`
		{
			"text": "Welcome",
			"alignment": "center",
			"justify": "center",
			"color": "#FF2244FF",
			"bg-color": "#002288FF"
		}
	`)

	// leagueId := "1259377273417502720" // "1236008483879403520"
	// leagueId :=  "1236008483879403520" // BH&TM
	// leagueId := "1253787285682409472"
	// t = "sleeper-matchups"
	// configJSON = []byte(`
	// 	{
	// 		"league_id": "` + leagueId + `",
	// 		"week": 1
	// 	}
	// `)

	// t = "nflbox"
	// configJSON = []byte(`{"auto": true}`)

	// go from []byte to specific ViewConfig type
	regView := viewCommon.RegisteredViews[t]
	configInstance := regView.NewConfig()
	if err := json.Unmarshal(configJSON, &configInstance); err != nil {
		log.Printf(fmt.Sprintf("Config for view type %s is invalid", t))
		return
	}

	newView, err := regView.NewView(configInstance)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create view of type %s with given config\nError: %s", t, err))
		return
	}

	// Set the initial view and start the application
	application.SetInitialView(newView)
	if err := application.Start(matrixDisplay, matrixConfig); err != nil {
		log.Printf("Failed to start application: %v", err)
		return
	}

	// Set the application instance for API handlers
	api.SetApplication(application)

	// Cleanup on exit
	defer application.Stop()

	// run the app server
	router := api.Router()
	router.Run(fmt.Sprintf("%s:%s", AppConfig.Server.AllowedHosts, AppConfig.Server.Port))

}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"log"

	d "github.com/6ixisgood/matrix-ticker/pkg/data"
	"github.com/6ixisgood/matrix-ticker/pkg/api"
	_ "github.com/6ixisgood/matrix-ticker/pkg/component"
	"github.com/6ixisgood/matrix-ticker/pkg/view"
	"github.com/6ixisgood/matrix-ticker/pkg/config"
	"github.com/6ixisgood/matrix-ticker/pkg/util"
	viewCommon "github.com/6ixisgood/matrix-ticker/pkg/view/common"
	"github.com/sixisgoood/go-rpi-rgb-led-matrix"

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
	configFilePath			= flag.String("config", "./config.yaml", "path to yaml config file")
)

func main() {
	flag.Parse()

	config.LoadConfig(configFilePath)	

	// set the RBG matrix configs
	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = config.AppConfig.Matrix.Rows
	matrixConfig.Cols = config.AppConfig.Matrix.Cols
	matrixConfig.Parallel = config.AppConfig.Matrix.Parallel
	matrixConfig.ChainLength = config.AppConfig.Matrix.Chain
	matrixConfig.Brightness = config.AppConfig.Matrix.Brightness
	matrixConfig.HardwareMapping = config.AppConfig.Matrix.HardwareMapping
	matrixConfig.ShowRefreshRate = config.AppConfig.Matrix.ShowRefresh
	matrixConfig.InverseColors = config.AppConfig.Matrix.InverseColors
	matrixConfig.DisableHardwarePulsing = config.AppConfig.Matrix.DisableHardwarePulsing
	matrixConfig.GpioSlowdown = config.AppConfig.Matrix.GpioSlowdown

	// configure the views
	viewCommon.SetViewConfig(&viewCommon.ViewConfig{
		MatrixRows: config.AppConfig.Matrix.Rows * config.AppConfig.Matrix.Parallel,
		MatrixCols: config.AppConfig.Matrix.Cols * config.AppConfig.Matrix.Chain,
		ImageDir: config.AppConfig.Data.ImageDir,
		CacheDir: config.AppConfig.Data.CacheDir,
		DefaultImageSizeX: config.AppConfig.Default.ImageSizeX,
		DefaultImageSizeY: config.AppConfig.Default.ImageSizeY,
		DefaultFontSize: config.AppConfig.Default.FontSize,
		DefaultFontColor: config.AppConfig.Default.FontColor,
		DefaultFontStyle: config.AppConfig.Default.FontStyle,
		DefaultFontType: config.AppConfig.Default.FontType,
	})

	// configure utils
	util.SetUtilConfig(&util.UtilConfig{
		CacheDir: config.AppConfig.Data.CacheDir,
		FontDir: config.AppConfig.Data.FontDir,
	})

	// configure server
	api.SetAppServerConfig(&api.AppServerConfig{
		AllowedHost: "127.0.0.1",
		Port: "8081",
	})

	// init the sports feed client
	d.InitSportsFeedClient(d.SportsFeedConfig{
		BaseUrl: config.AppConfig.Data.SportsFeed.BaseUrl,
		Username: config.AppConfig.Data.SportsFeed.Username,
		Password: config.AppConfig.Data.SportsFeed.Password,
	})

	// init the sleeper client
	d.InitSleeperClient(d.SleeperConfig{
		BaseUrl: config.AppConfig.Data.Sleeper.BaseUrl,
	})


	// setup matrix
	fmt.Println("Starting Matrix\n")
	m, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	fatal(err)

	tk := rgbmatrix.NewToolKit(m)
	defer tk.Close()

	// start the root animation
	animation := view.GetAnimation()

	log.Printf("Initializing the starting animation")
	args := map[string]string{
		"date": "20230910",
		"matchup": "20231207-NE-PIT",
		"src": "https://33.media.tumblr.com/ced5ea6f7722dd433465d2ab7e6e58e5/tumblr_nmt6p07KpV1ut1wfqo1_1280.gif",
		"league_id": "917441035486355456",
		"week": "3",
		"text": "Hello, World! It's raining outside right now",
		"layout": "flat",
	}
	animation.Init("nflbox", args)
	go tk.PlayAnimation(animation)

	// run the app server	
	api.Run()


}


func fatal(err error) {
	if err != nil {
		panic(err)
	}
}


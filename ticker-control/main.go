package main

import (
	"flag"
	"image"
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"	
	"time"

	d "github.com/sixisgoood/matrix-ticker/data"
	comp "github.com/sixisgoood/matrix-ticker/components"
	"github.com/sixisgoood/go-rpi-rgb-led-matrix"

)


var (
	configFilePath			= flag.String("config", "./config.yaml", "path to yaml config file")
	AppConfig				= ApplicationConfig{}
	Games					= d.DailyGamesNHLResponse{}
	Weather					= d.WeatherForecastResponse{}

)

var live_animation rgbmatrix.Animation

type ApplicationConfig struct {
	Matrix struct { 
		Rows					int		`yaml:"rows"`
		Cols					int		`yaml:"cols"`
		Parallel				int		`yaml:"parallel"`	
		Chain					int		`yaml:"chain"`
		Brightness				int		`yaml:"brightness"`
		HardwareMapping			string	`yaml:"harware_mapping"`
		ShowRefresh				bool	`yaml:"show_refresh"`
		InverseColors			bool	`yaml:"inverse_colors"`
		DisableHardwarePulsing	bool	`yaml:"disable_hardware_pulsing"`
	}	`yaml:"matrix"`
	Default struct {
		ImageSizeX				int		`yaml:"image_size_x"`
		ImageSizeY				int		`yaml:"image_size_y"`
		FontSize				int		`yaml:"font_size"`
		FontColor				string	`yaml:"font_color"`
		FontStyle				string	`yaml:"font_style"`
		FontType				string	`yaml:"font_type"`
	}
	Data struct {
		ImageDir				string	`yaml:"images"`
		CacheDir				string	`yaml:"cache"`
		SportsFeed struct {
			Username				string	`yaml:"username"`
			Password				string  `yaml:"password"`
		}	`yaml:"sportsfeed"`
		Weather struct {
			Key						string `yaml:"key"`
		}	`yaml:"weather"`
		
	}	`yaml:"data"`
}

type RootAnimation struct {}

func (a *RootAnimation) Next() (image.Image, <-chan time.Time, error){
	return live_animation.Next()
}

func getRootAnimation() (*RootAnimation) {
	return &RootAnimation{}
}

func SetLiveAnimation(new_animation rgbmatrix.Animation) {
	live_animation = new_animation
}

func GetApplicationConfig() ApplicationConfig {
	return AppConfig 
}

func GetGames() d.DailyGamesNHLResponse {
	return Games
}

func GetWeather() d.WeatherForecastResponse {
	return Weather
}

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatalf("Error loading config: '%v'", err)
		return
	}

	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		log.Fatalf("Error unmarshaling app config: '%v'", err)
		return
	}

	// configs
	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = AppConfig.Matrix.Rows
	matrixConfig.Cols = AppConfig.Matrix.Cols
	matrixConfig.Parallel = AppConfig.Matrix.Parallel
	matrixConfig.ChainLength = AppConfig.Matrix.Chain
	matrixConfig.Brightness = AppConfig.Matrix.Brightness
	matrixConfig.HardwareMapping = AppConfig.Matrix.HardwareMapping
	matrixConfig.ShowRefreshRate = AppConfig.Matrix.ShowRefresh
	matrixConfig.InverseColors = AppConfig.Matrix.InverseColors
	matrixConfig.DisableHardwarePulsing = AppConfig.Matrix.DisableHardwarePulsing

	// setup matrix
	fmt.Println("Starting Matrix\n")
	m, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	fatal(err)

	tk := rgbmatrix.NewToolKit(m)
	defer tk.Close()

	// start the root animation
	animation := GetAnimation()
	comp.SetViewGeneralConfig(comp.ViewGeneralConfig{
		MatrixRows: AppConfig.Matrix.Rows * AppConfig.Matrix.Parallel,
		MatrixCols: AppConfig.Matrix.Cols,
		ImageDir: AppConfig.Data.ImageDir,
		CacheDir: AppConfig.Data.CacheDir,
		SportsFeedUsername: AppConfig.Data.SportsFeed.Username,
		SportsFeedPassword: AppConfig.Data.SportsFeed.Password,
		DefaultImageSizeX: AppConfig.Default.ImageSizeX,
		DefaultImageSizeY: AppConfig.Default.ImageSizeY,
		DefaultFontSize: AppConfig.Default.FontSize,
		DefaultFontColor: AppConfig.Default.FontColor,
		DefaultFontStyle: AppConfig.Default.FontStyle,
		DefaultFontType: AppConfig.Default.FontType,
	})
	log.Printf("Initializing the starting animation")
	args := map[string]string{
		"date": "20230910",
		"src": "https://33.media.tumblr.com/ced5ea6f7722dd433465d2ab7e6e58e5/tumblr_nmt6p07KpV1ut1wfqo1_1280.gif",
	}
	animation.Init("image-player", args)
	go tk.PlayAnimation(animation)


	// start the http server
	server := NewAppServer()
	server.InitializeRoutes()
	server.Run("8081")

	// // start the server
	// Serve()	
}


func fatal(err error) {
	if err != nil {
		panic(err)
	}
}


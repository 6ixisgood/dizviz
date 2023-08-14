package main

import (
	"flag"
	"image"
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"	
	"time"

	//"github.com/sixisgoood/matrix-ticker/components"
	cd "github.com/sixisgoood/matrix-ticker/content_data"
	"github.com/sixisgoood/go-rpi-rgb-led-matrix"
)


var (
	configFilePath			= flag.String("config", "./config.yaml", "path to yaml config file")
	AppConfig				= ApplicationConfig{}
	Games					= cd.DailyGamesNHLResponse{}
	Weather					= cd.WeatherForecastResponse{}

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
	API	struct {
		NHL struct {
			Username				string	`yaml:"username"`
			Password				string  `yaml:"password"`
		}	`yaml:"nhl"`
		Weather struct {
			Key						string `yaml:"key"`
		}	`yaml:"weather"`
		
	}	`yaml:"api"`
	TemplateDir				string	`yaml:"template_dir`
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

func GetGames() cd.DailyGamesNHLResponse {
	return Games
}

func GetWeather() cd.WeatherForecastResponse {
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

	// config subpackages
	cd.NHLClientConfig = &cd.NHLRequestConfig{
		APIUsername: AppConfig.API.NHL.Username,
		APIPassword: AppConfig.API.NHL.Password,
	}
	cd.WeatherClientConfig = &cd.WeatherRequestConfig{
		Key: AppConfig.API.Weather.Key,
	}

	// setup matrix
	fmt.Println("Starting Matrix\n")
	m, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	fatal(err)

	tk := rgbmatrix.NewToolKit(m)
	defer tk.Close()

	// start the root animation
	animation := GetAnimation()
	log.Printf("Initializing the starting animation")
	args := map[string]string{
		"date": "20230110",
	}
	animation.Init("nhl-daily-games", args)
	tk.PlayAnimation(animation)


	// Games =  cd.FetchDailyNHLGamesInfo("2022-2023-regular", "20221012")
	// Weather = cd.FetchWeatherForecast("06105", "1")


	// // start the server
	// Serve()	
}


func fatal(err error) {
	if err != nil {
		panic(err)
	}
}


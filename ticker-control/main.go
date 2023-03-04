package main

import (
	"flag"
	"image"
	"fmt"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"	
	"time"

	"github.com/sixisgoood/matrix-ticker/animations"
	"github.com/sixisgoood/matrix-ticker/sports_data"
	"github.com/sixisgoood/go-rpi-rgb-led-matrix"
)


var (
	configFilePath			= flag.String("config", "./config.yaml", "path to yaml config file")
	rows                    = flag.Int("led-rows", 32, "number of rows supported")

)

var live_animation rgbmatrix.Animation

type AppConfig struct {
	Matrix struct { 
		Rows					int		`yaml:"rows"`
		Cols					int		`yaml:"cols"`
		Parallel				int		`yaml:"parallel"`	
		Chain					int		`yaml:"chain"`
		Brightness				int		`yaml:"brightness"`
		HardwareMapping			string	`yaml:"harware_mapping"`
		ShowRefresh				bool	`yaml:"show_refresh"`
		InverseColors			bool	`yaml:"inverse_colors"`
		DisableHardwarePulsing	bool	`yaml:"disable_hardware_pulsing`
	}	`yaml:"matrix"`
	API	struct {
		Username				string	`yaml:"username"`
		Password				string  `yaml:"password"`
	}	`yaml:"api"`
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

func main() {
	var appConfig AppConfig
	flag.Parse()

	data, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatalf("Error loading config: '%v'", err)
		return
	}

	if err := yaml.Unmarshal(data, &appConfig); err != nil {
		log.Fatalf("Error unmarshaling app config: '%v'", err)
		return
	}

	// configs
	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = appConfig.Matrix.Rows
	matrixConfig.Cols = appConfig.Matrix.Cols
	matrixConfig.Parallel = appConfig.Matrix.Parallel
	matrixConfig.ChainLength = appConfig.Matrix.Chain
	matrixConfig.Brightness = appConfig.Matrix.Brightness
	matrixConfig.HardwareMapping = appConfig.Matrix.HardwareMapping
	matrixConfig.ShowRefreshRate = appConfig.Matrix.ShowRefresh
	matrixConfig.InverseColors = appConfig.Matrix.InverseColors
	matrixConfig.DisableHardwarePulsing = appConfig.Matrix.DisableHardwarePulsing

	// config subpackages
	sports_data.ClientConfig = &sports_data.Config{
		APIUsername: appConfig.API.Username,
		APIPassword: appConfig.API.Password,
	}

	// setup matrix
	fmt.Println("Starting Matrix\n")
	m, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	fatal(err)

	tk := rgbmatrix.NewToolKit(m)
	defer tk.Close()

	// start the root animation
	content := `<matrix sizex="64" sizey="32"></matrix>`
	live_animation = animations.NewAnimation(content)
	animation := getRootAnimation()
	go tk.PlayAnimation(animation)

	// start the server
	Serve()	
}


func fatal(err error) {
	if err != nil {
		panic(err)
	}
}


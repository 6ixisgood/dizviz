package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
		Port string `yaml:"port"`
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

var (
	AppConfig = ApplicationConfig{}
)

func LoadConfig(filepath *string) {
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

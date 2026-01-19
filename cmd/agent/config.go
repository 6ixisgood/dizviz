package main

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

type AgentYAMLConfig struct {
	Agent struct {
		Name              string        `yaml:"name"`
		ControlPlaneAddr  string        `yaml:"control_plane_addr"`
		HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`
		ReconnectDelay    time.Duration `yaml:"reconnect_delay"`
		AutoReconnect     bool          `yaml:"auto_reconnect"`
		Version           string        `yaml:"version"`
	} `yaml:"agent"`

	Displays []struct {
		ID         string `yaml:"id"`
		Type       string `yaml:"type"` // "matrix", "hdmi", "file", etc.
		FPS        int    `yaml:"fps"`
		BufferSize int    `yaml:"buffer_size"`
		Hardware   struct {
			// Matrix-specific
			Rows                   int    `yaml:"rows"`
			Cols                   int    `yaml:"cols"`
			Parallel               int    `yaml:"parallel"`
			Chain                  int    `yaml:"chain"`
			Brightness             int    `yaml:"brightness"`
			HardwareMapping        string `yaml:"hardware_mapping"`
			ShowRefresh            bool   `yaml:"show_refresh"`
			InverseColors          bool   `yaml:"inverse_colors"`
			DisableHardwarePulsing bool   `yaml:"disable_hardware_pulsing"`
			GpioSlowdown           int    `yaml:"gpio_slowdown"`
			RateLimitHz            int    `yaml:"rate_limit_hz"`
		} `yaml:"hardware"`
	} `yaml:"displays"`

	Runtime struct {
		ImagesDir   string                            `yaml:"images_dir"`
		CacheDir    string                            `yaml:"cache_dir"`
		FontsDir    string                            `yaml:"fonts_dir"`
		DataSources map[string]map[string]interface{} `yaml:"data_sources"`
		Defaults    struct {
			FontSize  int    `yaml:"font_size"`
			FontColor string `yaml:"font_color"`
			FontStyle string `yaml:"font_style"`
			FontType  string `yaml:"font_type"`
		} `yaml:"defaults"`
	} `yaml:"runtime"`
}

var (
	Config = AgentYAMLConfig{}
)

func LoadConfig(filepath *string) {
	data, err := ioutil.ReadFile(*filepath)
	if err != nil {
		log.Fatalf("Error loading config: '%v'", err)
		return
	}

	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalf("Error unmarshaling agent config: '%v'", err)
		return
	}

	log.Printf("[Agent] Configuration loaded from %s", *filepath)
}

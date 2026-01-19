package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ControlPlaneConfig struct {
	Server struct {
		GRPCPort string `yaml:"grpc_port"`
		HTTPPort string `yaml:"http_port"`
	} `yaml:"server"`
	Data struct {
		StoreDir string `yaml:"store_dir"`
		Sleeper  struct {
			BaseUrl string `yaml:"base_url"`
		} `yaml:"sleeper"`
		SportsFeed struct {
			BaseUrl  string `yaml:"base_url"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"sportsfeed"`
		Weather struct {
			BaseUrl string `yaml:"base_url"`
			Key     string `yaml:"key"`
		} `yaml:"weather"`
	} `yaml:"data"`
}

var (
	Config = ControlPlaneConfig{}
)

func LoadConfig(filepath *string) {
	data, err := ioutil.ReadFile(*filepath)
	if err != nil {
		log.Fatalf("Error loading config: '%v'", err)
		return
	}

	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalf("Error unmarshaling control plane config: '%v'", err)
		return
	}

	log.Printf("[ControlPlane] Configuration loaded from %s", *filepath)
}

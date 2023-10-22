package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var configPath = "./config.yaml"

func GetConfig() Config {
	var cfg Config
	f := parseFlags()
	if f.ConfigPath != "" {
		configPath = f.ConfigPath
	}
	parseConfig(configPath, &cfg)

	if f.StoragePath != "" {
		cfg.StoragePath = f.StoragePath
	}
	if f.UploadPath != "" {
		cfg.UploadPath = f.UploadPath
	}
	if f.ServerAddress != "" {
		cfg.ServerAddress = f.ServerAddress
	}

	return cfg
}

type Config struct {
	StoragePath   string `yaml:"storage_path"`
	UploadPath    string `yaml:"upload_path"`
	ServerAddress string `yaml:"server_address"`
}

func parseConfig(filePath string, cfg any) {
	filename, _ := filepath.Abs(filePath)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		panic(err)
	}
}

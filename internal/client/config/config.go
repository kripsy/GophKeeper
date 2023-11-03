package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigPath    = "./config.yaml"
	defaultStoragePath   = "./1keeper/Data"
	defaultUploadPath    = "./1keeper/Upload"
	defaultServerAddress = "127.0.0.1:50051"
)

func GetConfig() Config {
	var cfg Config

	f := parseFlags()
	configPath := f.ConfigPath
	if configPath == "" {
		configPath = defaultConfigPath
	}

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		if err = parseConfig(configPath, &cfg); err != nil {
			fmt.Print("failed read yaml config file: ", err.Error())
			os.Exit(1)
		}
	}

	if f.StoragePath != "" {
		cfg.StoragePath = f.StoragePath
	}
	if f.UploadPath != "" {
		cfg.UploadPath = f.UploadPath
	}
	if f.ServerAddress != "" {
		cfg.ServerAddress = f.ServerAddress
	}

	if cfg.StoragePath == "" {
		cfg.StoragePath = defaultStoragePath
	}
	if cfg.UploadPath == "" {
		cfg.UploadPath = defaultUploadPath
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = defaultServerAddress
	}

	return cfg
}

type Config struct {
	StoragePath   string `yaml:"storage_path"`
	UploadPath    string `yaml:"upload_path"`
	ServerAddress string `yaml:"server_address"`
}

func parseConfig(filePath string, cfg any) error {
	filename, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return err
	}

	return nil
}

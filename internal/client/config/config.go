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
	var fileCfg Config

	flags := parseFlags()
	configPath := flags.ConfigPath
	if configPath == "" {
		configPath = defaultConfigPath
	}

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		if err = parseConfig(configPath, &fileCfg); err != nil {
			fmt.Print("failed read yaml config file: ", err.Error())
			os.Exit(1)
		}
	}

	return setConfig(fileCfg, flags)
}

func setConfig(fileCfg Config, flagCfg Flags) Config {
	if flagCfg.StoragePath != "" {
		fileCfg.StoragePath = flagCfg.StoragePath
	}
	if flagCfg.UploadPath != "" {
		fileCfg.UploadPath = flagCfg.UploadPath
	}
	if flagCfg.ServerAddress != "" {
		fileCfg.ServerAddress = flagCfg.ServerAddress
	}

	if fileCfg.StoragePath == "" {
		fileCfg.StoragePath = defaultStoragePath
	}
	if fileCfg.UploadPath == "" {
		fileCfg.UploadPath = defaultUploadPath
	}
	if fileCfg.ServerAddress == "" {
		fileCfg.ServerAddress = defaultServerAddress
	}

	return fileCfg
}

type Config struct {
	StoragePath   string `yaml:"storage_path"`
	UploadPath    string `yaml:"upload_path"`
	ServerAddress string `yaml:"server_address"`
}

func parseConfig(filePath string, cfg any) error {
	filename, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

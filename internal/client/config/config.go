package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigPath    = "./config.yaml"
	defaultStoragePath   = "./1keeper/Data"
	defaultUploadPath    = "./1keeper/Upload"
	defaultServerAddress = "127.0.0.1:50051"
)

var flags Flags
var once sync.Once

func GetConfig() (Config, error) {
	var fileCfg Config
	once.Do(func() {
		flags = parseFlags()
	})

	configPath := flags.ConfigPath
	if configPath == "" {
		configPath = defaultConfigPath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Возвращаем ошибку, если файл конфигурации не существует
		return Config{}, fmt.Errorf("config file not found: %s", configPath)
	} else if err != nil {
		// Возвращаем ошибку, если есть другая ошибка при проверке файла
		return Config{}, fmt.Errorf("failed to check config file: %w", err)
	}

	// Попытка парсинга файла конфигурации
	if err := parseConfig(configPath, &fileCfg); err != nil {
		fmt.Printf("failed to read yaml config file: %v\n", err)
		return Config{}, err
	}

	return setConfig(fileCfg, flags), nil
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

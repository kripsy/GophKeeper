package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Default values for configuration paths and server address.
const (
	defaultConfigPath    = "./config.yaml"
	defaultStoragePath   = "./1keeper/Data"
	defaultUploadPath    = "./1keeper/Upload"
	defaultServerAddress = "127.0.0.1:50051"
)

// Flags struct holds configuration flags.
//
//nolint:gochecknoglobals
var flags Flags

// Sync.Once ensures that configuration flags are parsed only once.
//
//nolint:gochecknoglobals
var once sync.Once

// GetConfig returns the configuration based on the provided flags or defaults.
func GetConfig() (Config, error) {
	var fileCfg Config
	once.Do(func() {
		flags = parseFlags()
	})

	// Use flag config path or the default if not provided.
	configPath := flags.ConfigPath
	if configPath == "" {
		configPath = defaultConfigPath
	}

	// Check if the specified config file exists.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Возвращаем ошибку, если файл конфигурации не существует
		//nolint:goerr113
		return Config{}, fmt.Errorf("config file not found: %s", configPath)
	} else if err != nil {
		// Возвращаем ошибку, если есть другая ошибка при проверке файла
		return Config{}, fmt.Errorf("failed to check config file: %w", err)
	}

	// Attempt to parse the configuration file.
	if err := parseConfig(configPath, &fileCfg); err != nil {
		fmt.Printf("failed to read yaml config file: %v\n", err)

		return Config{}, err
	}

	return setConfig(fileCfg, flags), nil
}

// setConfig updates the fileCfg with values from flagCfg if flags are provided.
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

// Config represents the configuration structure with YAML tags for unmarshaling.
type Config struct {
	StoragePath   string `yaml:"storage_path"`   // Path to store user data and secrets
	UploadPath    string `yaml:"upload_path"`    // Path to upload files from the secrets vault
	ServerAddress string `yaml:"server_address"` // Server address for sync
}

// parseConfig reads and parses the YAML configuration file.
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

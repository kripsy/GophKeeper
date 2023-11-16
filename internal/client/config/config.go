// Package config manages the configuration settings for the GophKeeper application.
// It provides functionality to parse and apply configuration from a YAML file
// and command-line flags.
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

// once ensures that configuration flags are parsed only once.
//
//nolint:gochecknoglobals
var once sync.Once

// GetConfig parses the configuration flags (if not already done) and then loads
// and applies the configuration settings from the YAML file or uses the defaults.
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
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		// Attempt to parse the configuration file.
		if err = parseConfig(configPath, &fileCfg); err != nil {
			fmt.Print("failed read yaml config file: ", err.Error())

			return Config{}, err
		}
	}

	return setConfig(fileCfg, flags), nil
}

// setConfig updates the fileCfg with values from flagCfg if flags are provided,
// otherwise, it uses default values.
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
// It includes settings for storage paths and server address.
type Config struct {
	StoragePath   string `yaml:"storage_path"`   // Path to store user data and secrets
	UploadPath    string `yaml:"upload_path"`    // Path to upload files from the secrets vault
	ServerAddress string `yaml:"server_address"` // Server address for sync
}

// parseConfig reads and parses the YAML configuration file, populating the cfg struct.
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

package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	validYAMLContent := []byte(`
storage_path: "./data"
upload_path: "./upload"
server_address: "127.0.0.1:50051"
`)
	invalidYAMLContent := []byte(`
storage_path: ["./data"]
`)
	//nolint:gosec
	require.NoError(t, os.WriteFile(configPath, validYAMLContent, 0644))

	tests := []struct {
		name        string
		filePath    string
		expectedCfg Config
		expectError bool
	}{
		{
			name:        "Valid config file",
			filePath:    configPath,
			expectedCfg: Config{StoragePath: "./data", UploadPath: "./upload", ServerAddress: "127.0.0.1:50051"},
			expectError: false,
		},
		{
			name:        "Invalid path",
			filePath:    string([]byte{0x7f}),
			expectedCfg: Config{},
			expectError: true,
		},
		{
			name:        "Non-existing file",
			filePath:    filepath.Join(tempDir, "non-existing.yaml"),
			expectedCfg: Config{},
			expectError: true,
		},
		{
			name:        "Invalid YAML content",
			filePath:    "",
			expectedCfg: Config{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			err := parseConfig(tt.filePath, &cfg)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedCfg, cfg)
			}
		})
	}

	//nolint:gosec
	require.NoError(t, os.WriteFile(configPath, invalidYAMLContent, 0644))
}

func TestSetConfig(t *testing.T) {
	tests := []struct {
		name    string
		fileCfg Config
		flagCfg Flags
		wantCfg Config
	}{
		{
			name:    "all defaults",
			fileCfg: Config{},
			flagCfg: Flags{},
			wantCfg: Config{StoragePath: defaultStoragePath, UploadPath: defaultUploadPath, ServerAddress: defaultServerAddress},
		},
		{
			name:    "all flags provided",
			fileCfg: Config{},
			flagCfg: Flags{StoragePath: "/custom/storage", UploadPath: "/custom/upload", ServerAddress: "192.168.1.1:6000"},
			wantCfg: Config{StoragePath: "/custom/storage", UploadPath: "/custom/upload", ServerAddress: "192.168.1.1:6000"},
		},
		{
			name:    "some flags provided",
			fileCfg: Config{},
			flagCfg: Flags{StoragePath: "/custom/storage"},
			wantCfg: Config{StoragePath: "/custom/storage", UploadPath: defaultUploadPath, ServerAddress: defaultServerAddress},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg := setConfig(tt.fileCfg, tt.flagCfg)
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("%s: setConfig() got = %v, want %v", tt.name, gotCfg, tt.wantCfg)
			}
		})
	}
}

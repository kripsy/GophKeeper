package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	// Предыдущая часть с созданием временного файла остается без изменений
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Восстанавливаем исходное значение после теста

	// Создаем временный файл конфигурации
	tempConfigFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Cannot create temp config file: %v", err)
	}
	defer os.Remove(tempConfigFile.Name()) // Удаляем временный файл после теста

	// Записываем валидный YAML в файл
	yamlContent := []byte(`
storage_path: "./1keeper/Data"
upload_path: "./1keeper/Upload"
server_address: "127.0.0.1:50051"
`)
	if _, err := tempConfigFile.Write(yamlContent); err != nil {
		t.Fatalf("Cannot write to temp config file: %v", err)
	}
	tempConfigFile.Close() // Закрываем файл после записи

	tests := []struct {
		name           string
		configFile     string // Добавляем поле для конфигурации пути файла
		configFileData []byte // Данные для записи в файл конфигурации
		args           []string
		want           config.Config
		wantErr        bool
	}{
		{
			name:       "ok with valid yaml",
			configFile: tempConfigFile.Name(),
			configFileData: []byte(`
storage_path: "./path/to/storage"
upload_path: "./path/to/uploads"
server_address: "127.0.0.1:8080"
`),
			args: append([]string{os.Args[0]},
				"-cfg-path", tempConfigFile.Name()),
			want: config.Config{
				StoragePath:   "./path/to/storage",
				UploadPath:    "./path/to/uploads",
				ServerAddress: "127.0.0.1:8080",
			},
			wantErr: false,
		},
		{
			name:           "error with invalid yaml",
			configFile:     tempConfigFile.Name(),
			configFileData: []byte(`invalid_yaml`),
			args:           append([]string{os.Args[0]}, "-cfg-path", tempConfigFile.Name()),
			want:           config.Config{},
			wantErr:        true,
		},
		{
			name:       "error with non-existing config path",
			configFile: "non-existing-path.yaml",
			configFileData: []byte(`
storage_path: "./path/to/storage"
upload_path: "./path/to/uploads"
server_address: "127.0.0.1:8080"
`),
			args:    append([]string{os.Args[0]}, "-cfg-path", "non-existing-path.yaml"),
			want:    config.Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Если указаны данные конфигурации, записываем их в файл
			if len(tt.configFileData) > 0 {
				//nolint:gosec
				if err := os.WriteFile(tt.configFile, tt.configFileData, 0644); err != nil {
					t.Fatalf("Failed to write to config file: %v", err)
				}
			}

			os.Args = tt.args
			got, err := config.GetConfig()

			if !tt.wantErr {
				require.Equal(t, true, reflect.DeepEqual(tt.want, got))
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			// Очищаем файл конфигурации после каждого теста
			if len(tt.configFileData) > 0 {
				os.Remove(tt.configFile)
			}
		})
	}
}

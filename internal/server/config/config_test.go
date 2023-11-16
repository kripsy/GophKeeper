package config_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/kripsy/GophKeeper/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	originalEnv := saveCurrentEnv()
	defer restoreEnv(originalEnv)
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *config.Config
	}{
		{
			name:    "Test Case 1 - Default Config",
			envVars: map[string]string{},
			expectedConfig: &config.Config{
				URLServer:            "localhost:8080",
				LoggerLevel:          "Info",
				DatabaseDsn:          "postgres://gophkeeperdb:gophkeeperdbpwd@localhost:5432/gophkeeperdb?sslmode=disable",
				Secret:               "supersecret",
				TokenExp:             config.TOKENEXP,
				IsSecure:             true,
				EndpointMinio:        "localhost:9000",
				AccessKeyIDMinio:     "masoud",
				SecretAccessKeyMinio: "Strong#Pass#2022",
				BucketNameMinio:      "secrets",
				IsUseSSLMinio:        false,
			},
		},
		{
			name: "Test Case 2 - Environment Variables",
			envVars: map[string]string{
				"SERVER_ADDRESS": "localhost:9090",
				"LOG_LEVEL":      "Debug",
				"SECRET":         "newsecret",
				"ISSECURE":       "false",
			},
			expectedConfig: &config.Config{
				URLServer:            "localhost:9090",
				LoggerLevel:          "Debug",
				DatabaseDsn:          "postgres://gophkeeperdb:gophkeeperdbpwd@localhost:5432/gophkeeperdb?sslmode=disable",
				Secret:               "newsecret",
				TokenExp:             config.TOKENEXP,
				IsSecure:             true,
				EndpointMinio:        "localhost:9000",
				AccessKeyIDMinio:     "masoud",
				SecretAccessKeyMinio: "Strong#Pass#2022",
				BucketNameMinio:      "secrets",
				IsUseSSLMinio:        false,
			},
		},
		{
			name: "Test Case 3 - Minio and DatabaseDsn Environment Variables",
			envVars: map[string]string{
				"ENDPOINTMINIO":        "minio.example.com:9000",
				"ACCESSKEYIDMINIO":     "minioadmin",
				"SECRETACCESSKEYMINIO": "minioadminpassword",
				"BUCKETNAMEMINIO":      "newbucket",
				"ISUSESSLMINIO":        "true",
				"DATABASE_DSN":         "postgres://user:password@localhost:5432/mydb?sslmode=enable",
			},
			expectedConfig: &config.Config{
				URLServer:            "localhost:8080",
				LoggerLevel:          "Info",
				DatabaseDsn:          "postgres://user:password@localhost:5432/mydb?sslmode=enable",
				Secret:               "supersecret",
				TokenExp:             config.TOKENEXP,
				IsSecure:             true,
				EndpointMinio:        "minio.example.com:9000",
				AccessKeyIDMinio:     "minioadmin",
				SecretAccessKeyMinio: "minioadminpassword",
				BucketNameMinio:      "newbucket",
				IsUseSSLMinio:        true,
			},
		},
	}

	tests[1].expectedConfig.IsSecure = false
	tests[2].expectedConfig.IsSecure = false
	tests[2].expectedConfig.URLServer = "localhost:9090"
	tests[2].expectedConfig.LoggerLevel = "Debug"
	tests[2].expectedConfig.Secret = "newsecret"
	tests[2].expectedConfig.DatabaseDsn = "postgres://user:password@localhost:5432/mydb?sslmode=enable"
	tests[2].expectedConfig.EndpointMinio = "minio.example.com:9000"
	tests[2].expectedConfig.AccessKeyIDMinio = "minioadmin"
	tests[2].expectedConfig.SecretAccessKeyMinio = "minioadminpassword"
	tests[2].expectedConfig.BucketNameMinio = "newbucket"
	tests[2].expectedConfig.IsUseSSLMinio = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer unsetEnvVars(tt.envVars)

			got, err := config.InitConfig()

			require.NoError(t, err)
			if !reflect.DeepEqual(tt.expectedConfig, got) {
				fmt.Println(tt.expectedConfig)
				fmt.Println(got)
			}
			assert.True(t, reflect.DeepEqual(tt.expectedConfig, got), "Configs do not match")
		})
	}
}

func setEnvVars(envVars map[string]string) {
	for key, value := range envVars {
		os.Setenv(key, value)
	}
}

func unsetEnvVars(envVars map[string]string) {
	for key := range envVars {
		os.Unsetenv(key)
	}
}

func saveCurrentEnv() map[string]string {
	env := os.Environ()
	savedEnv := make(map[string]string, len(env))
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		savedEnv[pair[0]] = pair[1]
	}

	return savedEnv
}

func restoreEnv(env map[string]string) {
	unsetEnvVars(env) // Удаляем текущие значения
	setEnvVars(env)   // Устанавливаем сохраненные значения
}

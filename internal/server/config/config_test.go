package config_test

import (
	"os"
	"testing"

	"github.com/kripsy/GophKeeper/internal/server/config"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
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
				IsSecure:             false,
				EndpointMinio:        "localhost:9000",
				AccessKeyIDMinio:     "masoud",
				SecretAccessKeyMinio: "Strong#Pass#2022",
				BucketNameMinio:      "secrets",
				IsUseSSLMinio:        false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			got, err := config.InitConfig()

			assert.NoError(t, err)

			assert.Equal(t, tt.expectedConfig, got)

			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

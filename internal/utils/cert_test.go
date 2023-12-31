package utils_test

import (
	"os"
	"testing"

	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateCertificate(t *testing.T) {
	tests := []struct {
		name           string
		serverCertPath string
		privateKeyPath string
		wantErr        bool
	}{
		{
			name:           "Valid paths",
			serverCertPath: "/tmp/server.crt",
			privateKeyPath: "/tmp/server.key",
			wantErr:        false,
		},
		{
			name:           "InValid path 1",
			serverCertPath: "/tmp/server.crt",
			privateKeyPath: " ",
			wantErr:        true,
		},
		{
			name:           "InValid path 2",
			serverCertPath: " ",
			privateKeyPath: "/tmp/server.key",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "cert_test")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)
			switch tt.name {
			case "Valid paths":
				tt.serverCertPath = tempDir + "/server.crt"
				tt.privateKeyPath = tempDir + "/server.key"
			case "InValid path 1":
				tt.serverCertPath = tempDir //+ "/server.crt"
				tt.privateKeyPath = tempDir + "/server.key"
			case "InValid path 2":
				tt.serverCertPath = tempDir + "/server.crt"
				tt.privateKeyPath = tempDir //+ "/server.key"
			}

			err = utils.CreateCertificate(tt.serverCertPath, tt.privateKeyPath)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if !tt.wantErr {
				_, err = os.Stat(tt.serverCertPath)
				require.NoError(t, err, "serverCertPath should exist after CreateCertificate")

				_, err = os.Stat(tt.privateKeyPath)
				require.NoError(t, err, "privateKeyPath should exist after CreateCertificate")
			}
		})
	}
}

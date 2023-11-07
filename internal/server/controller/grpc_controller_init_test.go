package controller_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/GophKeeper/internal/server/controller"
	"github.com/kripsy/GophKeeper/internal/server/controller/mocks"
	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInitGrpcServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	mockSecretUseCase := mocks.NewMockSecretUseCase(ctrl)

	logger := zap.NewNop()

	tests := []struct {
		name           string
		isSecure       bool
		serverCertPath string
		privateKeyPath string
		wantErr        bool
	}{
		{
			name:           "Secure server initialization",
			isSecure:       true,
			serverCertPath: "path/to/server.crt",
			privateKeyPath: "path/to/server.key",
			wantErr:        false,
		},
		{
			name:           "Insecure server initialization",
			isSecure:       false,
			serverCertPath: "",
			privateKeyPath: "",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isSecure {
				tempDir, err := os.MkdirTemp("", "cert_test")
				require.NoError(t, err)
				defer os.RemoveAll(tempDir)

				tt.serverCertPath = tempDir + "/server.crt"
				tt.privateKeyPath = tempDir + "/server.key"
				//nolint:errcheck
				utils.CreateCertificate(tt.serverCertPath, tt.privateKeyPath)
			}
			_, err := controller.InitGrpcServer(mockUserUseCase,
				mockSecretUseCase,
				"secret",
				tt.isSecure,
				tt.serverCertPath,
				tt.privateKeyPath, logger)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

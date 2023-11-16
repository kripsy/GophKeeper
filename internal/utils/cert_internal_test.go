package utils

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveCert(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "Valid data",
			data:    "test certificate data",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			payload := bytes.NewBufferString(tt.data)

			tempFile, err := os.CreateTemp("", "cert_test")
			require.NoError(err, "Failed to create temp file")
			defer os.Remove(tempFile.Name())

			err = saveCert(tempFile.Name(), payload)
			if tt.wantErr {
				require.Error(err, "saveCert() should return an error")
			} else {
				require.NoError(err, "saveCert() should not return an error")

				result, err := os.ReadFile(tempFile.Name())
				require.NoError(err, "Failed to read from temp file")
				require.Equal(tt.data, string(result), "The written data should match the test data")
			}
		})
	}
}

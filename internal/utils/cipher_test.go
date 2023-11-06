package utils

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {

	validKey := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, validKey)
	require.NoError(t, err, "Failed to generate a valid key")

	tests := []struct {
		name      string
		data      []byte
		cipherKey []byte
		wantErr   bool
	}{
		{
			name:      "ValidKey",
			data:      []byte("test data"),
			cipherKey: validKey,
			wantErr:   false,
		},
		{
			name:      "InvalidKeyShort",
			data:      []byte("test data"),
			cipherKey: []byte("short"),
			wantErr:   true,
		},
		{
			name:      "InvalidKeyLong",
			data:      []byte("test data"),
			cipherKey: make([]byte, 64),
			wantErr:   true,
		},
		{
			name:      "EmptyData",
			data:      []byte(""),
			cipherKey: validKey,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptedData, err := Encrypt(tt.data, tt.cipherKey)

			if tt.wantErr {
				assert.Error(t, err, "Encrypt() should return an error")
			} else {
				assert.NoError(t, err, "Encrypt() should not return an error")
				assert.NotEqual(t, tt.data, encryptedData, "Encrypted data should not match the original data")
			}
		})
	}
}

func TestDecrypt(t *testing.T) {

	cipherKey := make([]byte, 32)
	_, err := rand.Read(cipherKey)
	if err != nil {
		t.Fatalf("Failed to generate cipher key: %v", err)
	}

	originalText := "The quick brown fox jumps over the lazy dog"

	encryptedData, err := Encrypt([]byte(originalText), cipherKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	tests := []struct {
		name      string
		data      []byte
		cipherKey []byte
		wantErr   bool
	}{
		{
			name:      "Valid decryption",
			data:      encryptedData,
			cipherKey: cipherKey,
			wantErr:   false,
		},
		{
			name:      "Invalid key",
			data:      encryptedData,
			cipherKey: make([]byte, 32),
			wantErr:   true,
		},
		{
			name:      "Corrupted data",
			data:      append(encryptedData, byte(0)),
			cipherKey: cipherKey,
			wantErr:   true,
		},
		{
			name:      "Short data",
			data:      encryptedData[:10],
			cipherKey: cipherKey,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			decryptedData, err := Decrypt(tt.data, tt.cipherKey)

			if tt.wantErr {
				assert.Error(err, "Decrypt() should return an error")
			} else {
				assert.NoError(err, "Decrypt() should not return an error")
				assert.Equal(originalText, string(decryptedData), "The decrypted data should match the original")
			}
		})
	}
}

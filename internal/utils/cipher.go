package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// ErrCipherTextShort is an error returned when the provided ciphertext for decryption is shorter than expected.
var ErrCipherTextShort = errors.New("ciphertext too short")

// Encrypt encrypts the provided data using AES-GCM with the given cipher key.
//
// Parameters:
// - data: The byte slice of data to be encrypted.
// - cipherKey: The byte slice representing the key used for encryption.
//
// Returns:
// - A byte slice of the encrypted data.
// - An error if any issues occur during encryption.
func Encrypt(data []byte, cipherKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// Decrypt decrypts the provided data using AES-GCM with the given cipher key.
//
// Parameters:
// - data: The byte slice of data to be decrypted.
// - cipherKey: The byte slice representing the key used for decryption.
//
// Returns:
// - A byte slice of the decrypted data.
// - An error if any issues occur during decryption, including if the ciphertext is too short.
func Decrypt(data []byte, cipherKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("%w", ErrCipherTextShort)
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return plaintext, nil
}

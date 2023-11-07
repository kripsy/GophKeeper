package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

var ErrCipherTextShort = errors.New("ciphertext too short")

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

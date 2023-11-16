package models

import (
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/kripsy/GophKeeper/internal/utils"
)

// user metafile postfix.
const (
	metaPostfix = ".meta"
)

// UserData represents the user data structure, including User information and UserMeta.
type UserData struct {
	User User     `json:"user"`
	Meta UserMeta `json:"meta"`
}

// User represents user information, including username, password, and a key.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      []byte `json:"_,omitempty"`
}

// GetUserKey derives a cryptographic key based on the user's password and username.
func (u User) GetUserKey() ([]byte, error) {
	//nolint:wrapcheck
	return utils.DeriveKey(u.Password, u.Username)
}

// GetHashedPass derives a hashed password based on the user's username and password.
func (u User) GetHashedPass() (string, error) {
	hash, err := utils.DeriveKey(u.Username, u.Password)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return hex.EncodeToString(hash), err
}

// GetDir returns the directory path associated with the user's data.
func (u User) GetDir(dataPath string) string {
	return filepath.Join(dataPath, u.Username+metaPostfix)
}

// GetMetaFileName returns the meta file name associated with the user.
func (u User) GetMetaFileName() string {
	return u.Username + metaPostfix
}

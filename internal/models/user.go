package models

import (
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/kripsy/GophKeeper/internal/utils"
)

const (
	metaPostfix = ".meta"
)

type UserData struct {
	User User     `json:"user"`
	Meta UserMeta `json:"meta"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      []byte `json:"_,omitempty"`
}

func (u User) GetUserKey() ([]byte, error) {
	//nolint:wrapcheck
	return utils.DeriveKey(u.Password, u.Username)
}

func (u User) GetHashedPass() (string, error) {
	hash, err := utils.DeriveKey(u.Username, u.Password)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return hex.EncodeToString(hash), err
}

func (u User) GetDir(dataPath string) string {
	return filepath.Join(dataPath, u.Username+metaPostfix)
}

func (u User) GetMetaFileName() string {
	return u.Username + metaPostfix
}

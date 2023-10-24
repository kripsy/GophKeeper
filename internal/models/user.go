package models

import (
	"encoding/hex"
	"github.com/kripsy/GophKeeper/internal/utils"
	"path/filepath"
)

const (
	metaPostfix = ".meta"
)

type UserData struct {
	User User
	Meta UserMeta
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string
	Key      []byte `json:"_,omitempty"`
}

func (u User) GetUserKey() ([]byte, error) {
	return utils.DeriveKey(u.Password, u.Username)
}

func (u User) GetHashedPass() (string, error) {
	hash, err := utils.DeriveKey(u.Username, u.Password)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash), err
}

func (u User) GetDir(dataPath string) string {
	return filepath.Join(dataPath, u.Username+metaPostfix)
}

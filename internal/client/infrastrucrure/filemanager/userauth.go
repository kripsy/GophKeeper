package filemanager

import (
	"encoding/json"
	"errors"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"os"
)

const (
	fileMode os.FileMode = 0700 //660
	dirMode  os.FileMode = 0700 //755
)

type UserAuth struct {
	userFilePath string
}

func NewUserAuth(userPath string) (*UserAuth, error) {
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		if err = os.MkdirAll(userPath, dirMode); err != nil {
			return nil, err
		}
	}

	return &UserAuth{
		userFilePath: userPath,
	}, nil
}

func (a *UserAuth) IsUserNotExisting(userDit string) bool {
	if _, err := os.Stat(userDit); os.IsNotExist(err) {
		return true
	}

	return false
}

func (a *UserAuth) CreateUser(user *models.User, isLocalStorage bool) (models.UserMeta, error) {
	meta := models.UserMeta{
		Username:       user.Username,
		IsLocalStorage: isLocalStorage,
		Data:           make(map[string]models.DataInfo),
	}
	body, err := json.Marshal(meta)
	if err != nil {
		return models.UserMeta{}, err
	}

	key, err := user.GetUserKey()
	if err != nil {
		return models.UserMeta{}, err
	}

	encryptData, err := utils.Encrypt(body, key)
	if err != nil {
		return models.UserMeta{}, err
	}

	err = os.WriteFile(user.GetDir(a.userFilePath), encryptData, fileMode)
	if err != nil {
		return models.UserMeta{}, err
	}

	user.Key = key

	return meta, nil
}

func (a *UserAuth) GetUser(user *models.User) (models.UserMeta, error) {
	fileData, err := os.ReadFile(user.GetDir(a.userFilePath))
	if err != nil {
		return models.UserMeta{}, err
	}

	key, err := user.GetUserKey()
	if err != nil {
		return models.UserMeta{}, err
	}
	userData, err := utils.Decrypt(fileData, key)
	if err != nil {
		return models.UserMeta{}, err
	}

	meta := models.UserMeta{}
	err = json.Unmarshal(userData, &meta)
	if err != nil {
		return models.UserMeta{}, err
	}

	if user.Username != meta.Username {
		return models.UserMeta{}, errors.New("error compared user Data")
	}

	user.Key = key

	return meta, nil
}

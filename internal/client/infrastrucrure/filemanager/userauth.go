package filemanager

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
)

type Auth interface {
	IsUserNotExisting(userDit string) bool
	CreateUser(user *models.User, isLocalStorage bool) (models.UserMeta, error)
	GetUser(user *models.User) (models.UserMeta, error)
}

type userAuth struct {
	userFilePath string
}

//nolint:revive
func NewUserAuth(userPath string) (*userAuth, error) {
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		if err = os.MkdirAll(userPath, permissions.DirMode); err != nil {
			return nil, err
		}
	}

	return &userAuth{
		userFilePath: userPath,
	}, nil
}

func (a *userAuth) IsUserNotExisting(userDit string) bool {
	if _, err := os.Stat(userDit); os.IsNotExist(err) {
		return true
	}

	return false
}

func (a *userAuth) CreateUser(user *models.User, isLocalStorage bool) (models.UserMeta, error) {
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

	err = os.WriteFile(user.GetDir(a.userFilePath), encryptData, permissions.FileMode)
	if err != nil {
		return models.UserMeta{}, err
	}

	user.Key = key

	return meta, nil
}

func (a *userAuth) GetUser(user *models.User) (models.UserMeta, error) {
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

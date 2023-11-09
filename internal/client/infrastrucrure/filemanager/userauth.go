package filemanager

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils/crypto"
)

type Auth interface {
	IsUserNotExisting(userDit string) bool
	CreateUser(user *models.User, isSyncStorage bool) (models.UserMeta, error)
	GetUser(user *models.User) (models.UserMeta, error)
}

type userAuth struct {
	userFilePath string
}

//nolint:revive,nolintlint
func NewUserAuth(userPath string) (*userAuth, error) {
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		if err = os.MkdirAll(userPath, permissions.DirMode); err != nil {
			return nil, fmt.Errorf("%w", err)
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

func (a *userAuth) CreateUser(user *models.User, isSyncStorage bool) (models.UserMeta, error) {
	meta := models.UserMeta{
		Username:      user.Username,
		IsSyncStorage: isSyncStorage,
		Data:          make(models.MetaData),
		DeletedData:   make(models.Deleted),
	}
	body, err := json.Marshal(meta)
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}

	key, err := user.GetUserKey()
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}

	encryptData, err := crypto.Encrypt(body, key)
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}

	err = os.WriteFile(user.GetDir(a.userFilePath), encryptData, permissions.FileMode)
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}

	user.Key = key

	return meta, nil
}

func (a *userAuth) GetUser(user *models.User) (models.UserMeta, error) {
	fileData, err := os.ReadFile(user.GetDir(a.userFilePath))
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}

	key, err := user.GetUserKey()
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}
	userData, err := crypto.Decrypt(fileData, key)
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}

	meta := models.UserMeta{}
	err = json.Unmarshal(userData, &meta)
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}

	user.Key = key

	return meta, nil
}

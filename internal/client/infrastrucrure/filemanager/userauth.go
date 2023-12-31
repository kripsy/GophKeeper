package filemanager

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
)

// Auth interface defines methods for user local authentication.
type Auth interface {
	IsUseExisting(userDit string) bool
	CreateUser(user *models.User, isSyncStorage bool) (models.UserMeta, error)
	GetUser(user *models.User) (models.UserMeta, error)
}

// userAuth struct represents an authentication implementation using file storage.
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

// IsUserExisting checks if a user directory does not exist.
func (a *userAuth) IsUseExisting(userDit string) bool {
	if _, err := os.Stat(userDit); os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateUser creates a new user and stores their metadata in encrypted form.
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

	encryptData, err := utils.Encrypt(body, key)
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

// GetUser retrieves user metadata by decrypting it from the user's file.
func (a *userAuth) GetUser(user *models.User) (models.UserMeta, error) {
	fileData, err := os.ReadFile(user.GetDir(a.userFilePath))
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}

	key, err := user.GetUserKey()
	if err != nil {
		return models.UserMeta{}, fmt.Errorf("%w", err)
	}
	userData, err := utils.Decrypt(fileData, key)
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

package ui

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/models"
)

// UserInterface is the main interface for user interactions in the application,
// combining authentication, secret management, file handling, and synchronization.
type UserInterface interface {
	Auth
	Menu(isSyncStorage bool) int
	SecretManager
	FileDirector
	Sync(stop <-chan struct{})
	Clear()
	PrintErr(err string)
	Exit()
}

// Auth represents the authentication interface, handling user credentials and sync preferences.
type Auth interface {
	GetUser() (models.User, error)
	GetRepeatedPassword() (string, error)
	TryAgain() bool
	IsSyncStorage() bool
}

// SecretCreator defines methods for creating different types of secrets.
type SecretCreator interface {
	ChooseSecretType() (int, bool)
	AddNote() (filemanager.Note, error)
	AddBasicAuth() (filemanager.BasicAuth, error)
	AddCard() (filemanager.CardData, error)
	AddMetaInfo() (models.DataInfo, error)
}

// FileDirector handles file paths for writing and reading.
type FileDirector interface {
	UploadFileTo(cfgDir string) (string, bool)
	GetFilePath() string
}

// SecretManager combines the SecretCreator interface with methods for managing secrets.
type SecretManager interface {
	SecretCreator
	GetSecret(metaData models.MetaData) (string, bool)
	DeleteSecret(metaData models.MetaData) (string, bool)
	UpdateSecret(metaData models.MetaData) (string, int, bool)
}

package ui

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/models"
)

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

type Auth interface {
	GetUser() (models.User, error)
	GetRepeatedPassword() (string, error)
	TryAgain() bool
	IsLocalStorage() bool
}

type SecretCreator interface {
	ChooseSecretType() (int, bool)
	AddNote() (filemanager.Note, error)
	AddBasicAuth() (filemanager.BasicAuth, error)
	AddCard() (filemanager.CardData, error)
	AddMetaInfo() (models.DataInfo, error)
}

type FileDirector interface {
	UploadFileTo(cfgDir string) (string, bool)
	GetFilePath() string
}

type SecretManager interface {
	SecretCreator
	GetSecret(metaData models.MetaData) (string, bool)
	DeleteSecret(metaData models.MetaData) (string, bool)
	UpdateSecret(metaData models.MetaData) (string, int, bool)
}

package ui

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/models"
)

type UserInterface interface {
	GetUser() (models.User, error)
	GetRepeatedPassword() (string, error)
	TryAgain() bool
	IsLocalStorage() bool

	Menu(isLocalStorage bool) int

	GetSecret(metaData models.MetaData) (string, bool)
	DeleteSecret(metaData models.MetaData) (string, bool)
	UpdateSecret(metaData models.MetaData) (string, int, bool)

	ChooseSecretType() (int, bool)
	AddNote() (filemanager.Note, error)
	AddBasicAuth() (filemanager.BasicAuth, error)
	AddCard() (filemanager.CardData, error)
	AddMetaInfo() (models.DataInfo, error)

	UploadFileTo(cfgDir string) (string, bool)
	GetFilePath() string

	Sync(stop <-chan struct{})
	Clear()
	PrintErr(error string)
	Exit()
}

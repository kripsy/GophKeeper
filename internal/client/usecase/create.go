package usecase

import (
	"os"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
)

func (c *ClientUsecase) createSecret(secretType int, success bool) {
	defer c.InMenu()
	if !success {
		c.ui.PrintErr(ui.CreateErr)

		return
	}

	var data filemanager.Data
	var info models.DataInfo
	var err error

	data, info, err = c.getUserData(secretType)
	if err != nil {
		c.ui.PrintErr(ui.CreateErr)
		c.log.Err(err).Msg("failed to get user data")
		return
	}

	info.DataType = secretType
	err = c.fileManager.AddToStorage(info.Name, data, info)
	if err != nil {
		c.ui.PrintErr(ui.CreateErr)
		c.log.Err(err).Msg(ui.CreateErr)
		return
	}
}

//nolint:nolintlint
func (c *ClientUsecase) getUserData(secretType int) (filemanager.Data, models.DataInfo, error) {
	var data filemanager.Data
	var info models.DataInfo
	var err error

	switch secretType {
	case filemanager.NoteType:
		data, err = c.ui.AddNote()
	case filemanager.BasicAuthType:
		data, err = c.ui.AddBasicAuth()
	case filemanager.CardDataType:
		data, err = c.ui.AddCard()
	case filemanager.FileType:
		path := c.ui.GetFilePath()
		info.SetFileName(path)
		body, err := os.ReadFile(path)
		if err == nil {
			data = filemanager.File{Data: body}
		}
	}

	if err == nil {
		info, err = c.ui.AddMetaInfo()
	}

	return data, info, err
}

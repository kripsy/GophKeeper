//nolint:ireturn
package usecase

import (
	"fmt"
	"os"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
)

func (c *ClientUsecase) createSecret(secretType int, success bool) {
	if !success {
		return
	}

	var data filemanager.Data
	var info models.DataInfo
	var err error

	data, info, err = c.getSecretrData(secretType)
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

func (c *ClientUsecase) getSecretrData(secretType int) (filemanager.Data, models.DataInfo, error) {
	var data filemanager.Data
	var filePath string
	var err error

	switch secretType {
	case filemanager.NoteType:
		data, err = c.ui.AddNote()
	case filemanager.BasicAuthType:
		data, err = c.ui.AddBasicAuth()
	case filemanager.CardDataType:
		data, err = c.ui.AddCard()
	case filemanager.FileType:
		filePath = c.ui.GetFilePath()
		var body []byte
		body, err = os.ReadFile(filePath)
		data = filemanager.File{Data: body}
	}
	if err != nil {
		return nil, models.DataInfo{}, fmt.Errorf("%w", err)
	}
	info, err := c.ui.AddMetaInfo()
	if err != nil {
		return nil, models.DataInfo{}, fmt.Errorf("%w", err)
	}

	if secretType == filemanager.FileType {
		info.SetFileName(filePath)
	}

	return data, info, nil
}

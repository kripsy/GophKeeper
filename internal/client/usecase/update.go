//nolint:ireturn
package usecase

import (
	"fmt"
	"os"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
)

func (c *ClientUsecase) updateSecret(secretName string, updateType int, success bool) {
	defer c.InMenu()
	if !success {
		return
	}

	metaInfo := c.userData.Meta.Data[secretName]

	if ui.UpdateItems[updateType] == ui.Info {
		if err := c.updateMetaInfo(secretName, metaInfo); err != nil {
			c.ui.PrintErr(ui.UpdateErr)
			c.log.Err(err).Msg(ui.UpdateErr)
		}

		return
	}

	data, err := c.getUpdatedData(secretName, metaInfo.DataType)
	if err != nil {
		c.ui.PrintErr(ui.UpdateErr)
		c.log.Err(err).Msg(ui.UpdateErr)

		return
	}

	if err = c.fileManager.UpdateDataByName(secretName, data); err != nil {
		c.ui.PrintErr(ui.UpdateErr)
		c.log.Err(err).Msg(ui.UpdateErr)
	}
}

func (c *ClientUsecase) updateMetaInfo(secretName string, metaInfo models.DataInfo) error {
	fmt.Printf("Before  「 Name: %s, Description: %s 」\n", metaInfo.Name, metaInfo.Description)
	info, err := c.ui.AddMetaInfo()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = c.fileManager.UpdateInfoByName(secretName, info); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *ClientUsecase) getUpdatedData(secretName string, dataType int) (filemanager.Data, error) {
	var data filemanager.Data
	var err error

	switch dataType {
	case filemanager.NoteType:
		data, err = c.ui.AddNote()
	case filemanager.BasicAuthType:
		data, err = c.ui.AddBasicAuth()
	case filemanager.CardDataType:
		data, err = c.ui.AddCard()
	case filemanager.FileType:
		path := c.ui.GetFilePath()
		newInfo := models.DataInfo{}
		newInfo.SetFileName(path)
		body, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		err = c.fileManager.UpdateInfoByName(secretName, newInfo)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		data = filemanager.File{Data: body}
	}

	return data, fmt.Errorf("%w", err)
}

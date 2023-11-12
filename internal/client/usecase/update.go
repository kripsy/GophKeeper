package usecase

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
)

// updateSecret updates a secret based on the provided parameters.
func (c *ClientUsecase) updateSecret(secretName string, updateType int, success bool) {
	if !success {
		return
	}

	metaInfo := c.userData.Meta.Data[secretName]

	// Choosing the type of upgrade
	if ui.UpdateItems[updateType] == ui.Info {
		if err := c.updateMetaInfo(secretName, metaInfo); err != nil {
			c.ui.PrintErr(ui.UpdateErr)
			c.log.Err(err).Msg(ui.UpdateErr)
		}

		return
	}

	if metaInfo.DataType == filemanager.FileType {
		path := c.ui.GetFilePath()
		newInfo := metaInfo
		newInfo.SetFileName(path)
		err := c.fileManager.AddFileToStorage(false, newInfo.Name, path, newInfo)
		if err != nil {
			c.ui.PrintErr(ui.UpdateErr)
			c.log.Err(err).Msg(ui.UpdateErr)

			return
		}
		err = c.fileManager.UpdateInfoByName(secretName, newInfo)
		if err != nil {
			c.ui.PrintErr(ui.UpdateErr)
			c.log.Err(err).Msg(ui.UpdateErr)

			return
		}

		return
	}

	data, err := c.getUpdatedData(metaInfo.DataType)
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

// updateMetaInfo updates metadata information for a secret.
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

// getUpdatedData retrieves updated data based on the secret's type.
//
//nolint:ireturn,nolintlint
func (c *ClientUsecase) getUpdatedData(dataType int) (filemanager.Data, error) {
	var data filemanager.Data
	var err error

	// Handle different secret types.
	switch dataType {
	case filemanager.NoteType:
		data, err = c.ui.AddNote()
	case filemanager.BasicAuthType:
		data, err = c.ui.AddBasicAuth()
	case filemanager.CardDataType:
		data, err = c.ui.AddCard()
	}
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

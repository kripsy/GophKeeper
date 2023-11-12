package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
)

func (c *ClientUsecase) createSecret(secretType int, success bool) {
	if !success {
		return
	}

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
	}
	if err != nil {
		c.ui.PrintErr(ui.CreateErr)
		c.log.Err(err).Msg(ui.CreateErr)

		return
	}

	if secretType == filemanager.FileType {
		if err = c.addFileToStorage(filePath); err != nil {
			c.ui.PrintErr(ui.CreateErr)
			c.log.Err(err).Msg(ui.CreateErr)
		}

		return
	}

	if err = c.addToStorage(secretType, data); err != nil {
		c.ui.PrintErr(ui.CreateErr)
		c.log.Err(err).Msg(ui.CreateErr)
	}
}

func (c *ClientUsecase) addToStorage(secretType int, data filemanager.Data) error {
	info, err := c.ui.AddMetaInfo()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	info.DataType = secretType
	err = c.fileManager.AddToStorage(info.Name, data, info)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *ClientUsecase) addFileToStorage(filePath string) error {
	info, err := c.ui.AddMetaInfo()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	info.SetFileName(filePath)
	info.DataType = filemanager.FileType
	info.DataID = uuid.New().String()

	err = c.fileManager.AddFileToStorage(true, info.Name, filePath, info)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

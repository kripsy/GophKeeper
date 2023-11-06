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
	var err error

	if ui.UpdateItems[updateType] == ui.Info {
		err = c.updateMetaInfo(secretName, metaInfo)
	} else {
		data, err := c.getUpdatedData(secretName, metaInfo.DataType)
		if err == nil {
			err = c.fileManager.UpdateDataByName(secretName, data)
		}
	}

	if err != nil {
		c.ui.PrintErr(ui.UpdateErr)
		c.log.Err(err).Msg(ui.UpdateErr)
	}
}

func (c *ClientUsecase) updateMetaInfo(secretName string, metaInfo models.DataInfo) error {
	fmt.Println(fmt.Sprintf("Before  「 Name: %s, Description: %s 」", metaInfo.Name, metaInfo.Description))
	info, err := c.ui.AddMetaInfo()
	if err != nil {
		return err
	}
	return c.fileManager.UpdateInfoByName(secretName, info)
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
		if err == nil {
			err = c.fileManager.UpdateInfoByName(secretName, newInfo)
			data = filemanager.File{Data: body}
		}
	}

	return data, err
}

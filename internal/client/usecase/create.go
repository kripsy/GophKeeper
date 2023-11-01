package usecase

import (
	"fmt"
	"os"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
)

func (c *ClientUsecase) createSecret(secretType int, success bool) {
	defer c.InMenu()
	if !success {
		return
	}

	//todo handle error
	switch secretType {
	case filemanager.NoteType:
		data, _ := c.ui.AddNote()
		info, _ := c.ui.AddMetaInfo()
		info.DataType = secretType
		c.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.BasicAuthType:
		data, _ := c.ui.AddBasicAuth()
		info, _ := c.ui.AddMetaInfo()
		info.DataType = secretType
		c.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.CardDataType:
		data, _ := c.ui.AddCard()
		info, _ := c.ui.AddMetaInfo()
		info.DataType = secretType
		c.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.FileType:
		path := c.ui.GetFilePath()
		info, _ := c.ui.AddMetaInfo()
		info.DataType = secretType
		info.SetFileName(path)
		fmt.Println(info)
		body, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		data := filemanager.File{body}
		c.fileManager.AddToStorage(info.Name, data, info)
	}
}

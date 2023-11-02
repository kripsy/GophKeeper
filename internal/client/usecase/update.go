package usecase

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
)

func (c *ClientUsecase) updateSecret(secretName string, updateType int, success bool) {
	defer c.InMenu()
	if !success {
		return
	}
	metaInfo := c.userData.Meta.Data[secretName]

	if ui.UpdateItems[updateType] == ui.Info {
		fmt.Println(fmt.Sprintf("Before  「 Name: %s, Description: %s 」", metaInfo.Name, metaInfo.Description))
		info, _ := c.ui.AddMetaInfo()
		c.fileManager.UpdateInfoByName(secretName, info)
		return
	}

	var data filemanager.Data
	switch metaInfo.DataType {
	case filemanager.NoteType:
		data, _ = c.ui.AddNote()
	case filemanager.BasicAuthType:
		data, _ = c.ui.AddBasicAuth()
	case filemanager.CardDataType:
		data, _ = c.ui.AddCard()
		//case filemanager.FileType:
		//	var path string        // todo нужно обновлять вместе с FileName что содержится в мета
		//	cli.GetFilePath(&path)
		//	info.SetFileName(path)
		//	body, err := os.ReadFile(path)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//	data := filemanager.File{body}
	}
	c.fileManager.UpdateDataByName(secretName, data)
}

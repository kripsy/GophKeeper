package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"os"
	"path/filepath"
)

func (c *ClientUsecase) getSecrets(secretName string, success bool) {
	defer c.InMenu()
	if !success {
		return
	}

	data, info, err := c.fileManager.GetByName(secretName)
	if err != nil {
		c.log.Err(err).Msg("err get secret from meta")
	}

	switch info.DataType {
	case filemanager.NoteType:
		var dataStruct filemanager.Note
		_ = json.Unmarshal(data, &dataStruct)
		fmt.Println(fmt.Sprintf("%s, %s: %v", info.Name, info.Description, dataStruct)) //todo odod

	case filemanager.BasicAuthType:
		var dataStruct filemanager.BasicAuth
		_ = json.Unmarshal(data, &dataStruct)
		fmt.Println(fmt.Sprintf("%s, %s: %v", info.Name, info.Description, dataStruct)) //todo odod

	case filemanager.CardDataType:
		var dataStruct filemanager.CardData
		_ = json.Unmarshal(data, &dataStruct)
		fmt.Println(fmt.Sprintf("%s, %s: %v", info.Name, info.Description, dataStruct)) //todo odod

	case filemanager.FileType:
		var dataStruct filemanager.File
		newFilePath, ok := c.ui.UploadFileTo(c.uploadPath)
		if !ok {
			return
		}
		_ = json.Unmarshal(data, &dataStruct)
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			if err = os.MkdirAll(newFilePath, 0777); err != nil {
				fmt.Println(err)
			}
		}
		err = os.WriteFile(filepath.Join(newFilePath, *info.FileName), dataStruct.Data, 0777)
		if err != nil {
			fmt.Println("Success upload")
		}
		fmt.Println(info)
		fmt.Println(err)
	}
}

package usecase

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/manifoldco/promptui"
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
		err = json.Unmarshal(data, &dataStruct)
		if err != nil {
			c.log.Err(err).Msg("failed unmarshal data")

			return
		}
		printSecret(info.Name, info.Description, dataStruct.String())

	case filemanager.BasicAuthType:
		var dataStruct filemanager.BasicAuth
		err = json.Unmarshal(data, &dataStruct)
		if err != nil {
			c.log.Err(err).Msg("failed unmarshal data")

			return
		}
		printSecret(info.Name, info.Description, dataStruct.String())

	case filemanager.CardDataType:
		var dataStruct filemanager.CardData
		err = json.Unmarshal(data, &dataStruct)
		if err != nil {
			c.log.Err(err).Msg("failed unmarshal data")

			return
		}
		printSecret(info.Name, info.Description, dataStruct.String())

	case filemanager.FileType:
		var dataStruct filemanager.File
		newFilePath, ok := c.ui.UploadFileTo(c.uploadPath)
		if !ok {
			return
		}
		err = json.Unmarshal(data, &dataStruct)
		if err != nil {
			c.log.Err(err).Msg("failed unmarshal data")

			return
		}
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			if err = os.MkdirAll(newFilePath, 0777); err != nil {
				fmt.Println(err)
			}
		}
		err = os.WriteFile(filepath.Join(newFilePath, *info.FileName), dataStruct.Data, 0777)
		if err != nil {
			c.log.Err(err).Msg("failed write secret file")

			return
		}

		fmt.Println(newFilePath)
		printSecret(info.Name, info.Description, dataStruct.String())
	}
}

func printSecret(name, description, secret string) {
	fmt.Println(promptui.Styler(
		promptui.FGBold,
		promptui.FGBlue,
	)(fmt.Sprintf(" Name: %s \n Description: %s \n %s", name, description, secret)))
}

package usecase

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

// getSecrets retrieves and displays information about a secret based on its name.
func (c *ClientUsecase) getSecrets(secretName string, success bool) {
	if !success {
		return
	}

	info, ok := c.userData.Meta.Data[secretName]
	if !ok {
		c.ui.PrintErr(ui.GetErr)
		c.log.Error().Msg("err get secret from meta")

		return
	}

	if info.DataType == filemanager.FileType {
		c.getFileSecret(info)

		return
	}

	// Retrieve data and information about the secret.
	data, info, err := c.fileManager.GetByInfo(info)
	if err != nil {
		c.ui.PrintErr(ui.GetErr)
		c.log.Err(err).Msg("err get secret from meta")

		return
	}
	// Choose the appropriate handling based on the data type of the secret.
	switch info.DataType {
	case filemanager.NoteType:
		var dataStruct filemanager.Note
		c.getSecret(data, info, &dataStruct)
	case filemanager.BasicAuthType:
		var dataStruct filemanager.BasicAuth
		c.getSecret(data, info, &dataStruct)
	case filemanager.CardDataType:
		var dataStruct filemanager.CardData
		c.getSecret(data, info, &dataStruct)
	}
}

// getSecret displays the information of a non-file type secret.
func (c *ClientUsecase) getSecret(data []byte, info models.DataInfo, dataStruct filemanager.Data) {
	err := json.Unmarshal(data, dataStruct)
	if err != nil {
		c.ui.PrintErr(ui.GetErr)
		c.log.Err(err).Msg("failed unmarshal data")

		return
	}
	printSecret(info.Name, info.Description, dataStruct.String())
}

// getFileSecret prompts you to choose where to upload the file and displays a success message.
func (c *ClientUsecase) getFileSecret(info models.DataInfo) {
	newFilePath, ok := c.ui.UploadFileTo(c.uploadPath)
	if !ok {
		c.ui.PrintErr(ui.GetErr)

		return
	}

	if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
		if err = os.MkdirAll(newFilePath, permissions.DirMode); err != nil {
			c.ui.PrintErr(ui.GetErr)

			c.log.Err(err).Msg("failed make dir for secret file")

			return
		}
	}

	err := c.fileManager.ReadFileFromStorage(newFilePath, info)
	if err != nil {
		c.ui.PrintErr(ui.GetErr)

		c.log.Err(err).Msg("failed write secret file")

		return
	}

	printSecret(info.Name, info.Description, newFilePath)
}

// printSecret prints formatted information about a secret.
func printSecret(name, description, secret string) {
	fmt.Println(promptui.Styler(
		promptui.FGBold,
		promptui.FGBlue,
	)(fmt.Sprintf(" Name: %s \n Description: %s \n %s", name, description, secret)))
}

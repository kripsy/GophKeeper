package app

import (
	"encoding/json"
	"fmt"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/cli"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"os"
	"path/filepath"
)

func (a *Application) inMenu() {
	action := cli.MenuTable[cli.Menu()]
	switch action {
	case cli.SecretsKey:
		a.getSecrets(cli.GetSecret(a.userData.Meta.Data)) // todo сторит обдумать метадату на name
	case cli.AddSecretKey:
		a.createSecret(cli.SecretType())
	case cli.UpdateSecretKey: //todo идут строительные работы
		a.updateSecret(cli.UpdateSecret(a.userData.Meta.Data))
	case cli.DeleteSecretKey:
		a.deleteSecret(cli.DeleteSecret(a.userData.Meta.Data))
	case cli.ExitKey:
		cli.Exit()
	}
}

func (a *Application) getSecrets(secretName string, success bool) {
	defer a.inMenu()
	if !success {
		return
	}

	data, info, err := a.fileManager.GetByName(secretName)

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
		newFilePath, ok := cli.UploadFileTo(a.uploadPath)
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
			fmt.Println("Sucsess upload")
		}
		fmt.Println(info)
		fmt.Println(err)
	}
}

func (a *Application) createSecret(secretType int, success bool) {
	defer a.inMenu()
	if !success {
		return
	}

	//todo handle error
	switch secretType {
	case filemanager.NoteType:
		data, _ := cli.AddNote()
		info, _ := cli.AddMetaInfo()
		info.DataType = secretType
		a.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.BasicAuthType:
		data, _ := cli.AddBasicAuth()
		info, _ := cli.AddMetaInfo()
		info.DataType = secretType
		a.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.CardDataType:
		data, _ := cli.AddCard()
		info, _ := cli.AddMetaInfo()
		info.DataType = secretType
		a.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.FileType:
		var path string
		cli.GetFilePath(&path)
		info, _ := cli.AddMetaInfo()
		info.DataType = secretType
		info.SetFileName(path)
		fmt.Println(info)
		body, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		data := filemanager.File{body}
		a.fileManager.AddToStorage(info.Name, data, info)
	}
}
func (a *Application) updateSecret(secretName string, updateType int, success bool) {
	defer a.inMenu()
	if !success {
		return
	}
	metaInfo := a.userData.Meta.Data[secretName]

	if cli.UpdateItems[updateType] == cli.Info {
		fmt.Println(fmt.Sprintf("Before  「 Name: %s, Description: %s 」", metaInfo.Name, metaInfo.Description))
		info, _ := cli.AddMetaInfo()
		a.fileManager.UpdateInfoByName(secretName, info)
		return
	}

	var data filemanager.Data
	switch metaInfo.DataType {
	case filemanager.NoteType:
		data, _ = cli.AddNote()
	case filemanager.BasicAuthType:
		data, _ = cli.AddBasicAuth()
	case filemanager.CardDataType:
		data, _ = cli.AddCard()
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
	a.fileManager.UpdateDataByName(secretName, data)
}

func (a *Application) deleteSecret(secretName string, success bool) {
	defer a.inMenu()
	if !success {
		return
	}

	a.fileManager.DeleteByName(secretName)
}

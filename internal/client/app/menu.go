package app

import (
	"encoding/json"
	"fmt"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/manifoldco/promptui"
	"os"
	"path/filepath"
)

func (a *Application) inMenu() {
	action := ui.MenuTable[a.ui.Menu(a.userData.Meta.IsLocalStorage)]
	switch action {
	case ui.SecretsKey:
		a.getSecrets(a.ui.GetSecret(a.userData.Meta.Data)) // todo сторит обдумать метадату на name
	case ui.AddSecretKey:
		a.createSecret(a.ui.ChooseSecretType())
	case ui.UpdateSecretKey: //todo идут строительные работы
		a.updateSecret(a.ui.UpdateSecret(a.userData.Meta.Data))
	case ui.SyncSecrets:
		a.sync()
	case ui.DeleteSecretKey:
		a.deleteSecret(a.ui.DeleteSecret(a.userData.Meta.Data))
	case ui.About:
		a.about()
	case ui.ExitKey:
		a.ui.Exit()
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
		newFilePath, ok := a.ui.UploadFileTo(a.uploadPath)
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

func (a *Application) createSecret(secretType int, success bool) {
	defer a.inMenu()
	if !success {
		return
	}

	//todo handle error
	switch secretType {
	case filemanager.NoteType:
		data, _ := a.ui.AddNote()
		info, _ := a.ui.AddMetaInfo()
		info.DataType = secretType
		a.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.BasicAuthType:
		data, _ := a.ui.AddBasicAuth()
		info, _ := a.ui.AddMetaInfo()
		info.DataType = secretType
		a.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.CardDataType:
		data, _ := a.ui.AddCard()
		info, _ := a.ui.AddMetaInfo()
		info.DataType = secretType
		a.fileManager.AddToStorage(info.Name, data, info)
	case filemanager.FileType:
		path := a.ui.GetFilePath()
		info, _ := a.ui.AddMetaInfo()
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

	if ui.UpdateItems[updateType] == ui.Info {
		fmt.Println(fmt.Sprintf("Before  「 Name: %s, Description: %s 」", metaInfo.Name, metaInfo.Description))
		info, _ := a.ui.AddMetaInfo()
		a.fileManager.UpdateInfoByName(secretName, info)
		return
	}

	var data filemanager.Data
	switch metaInfo.DataType {
	case filemanager.NoteType:
		data, _ = a.ui.AddNote()
	case filemanager.BasicAuthType:
		data, _ = a.ui.AddBasicAuth()
	case filemanager.CardDataType:
		data, _ = a.ui.AddCard()
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

func (a *Application) sync() {
	fmt.Println("не реализовано")
}

func (a *Application) deleteSecret(secretName string, success bool) {
	defer a.inMenu()
	if !success {
		return
	}

	a.fileManager.DeleteByName(secretName)
}

func (a *Application) about() {
	defer a.inMenu()
	fmt.Println(promptui.Styler(
		promptui.FGItalic,
		promptui.BGBlue,
	)(fmt.Sprintf(aboutMsg, a.buildInfo.BuildVersion, a.buildInfo.BuildDate)))
}

const aboutMsg = `
「  GophKeeper  」
This is an Application with the ability to store your secrets locally,
as well as synchronize between your clients when registering through the server.
Build version: %s                          Build date: %s
`

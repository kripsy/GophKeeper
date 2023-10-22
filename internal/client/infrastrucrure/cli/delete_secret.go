package cli

import (
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

func DeleteSecret(metaData models.MetaData) (string, bool) {
	name, ok := chooseSecret(metaData, DeleteSecretKey, deleteSecretTemplate)
	if !ok {
		return "", ok
	}

	isDelete := promptui.Prompt{
		Label:     "Secret will be deleted from the device, are you sure?",
		IsConfirm: true,
	}

	_, err := isDelete.Run()
	if err != nil {
		return "", false
	}

	return name, true
}

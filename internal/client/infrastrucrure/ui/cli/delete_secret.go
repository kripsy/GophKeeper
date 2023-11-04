package cli

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

func (c *CLI) DeleteSecret(metaData models.MetaData) (string, bool) {
	name, ok := c.chooseSecret(metaData, ui.DeleteSecretKey, deleteSecretTemplate)
	if !ok {
		return "", ok
	}

	isDelete := promptui.Prompt{
		Label:       "Secret will be deleted from the device, are you sure?",
		IsConfirm:   true,
		HideEntered: true,
	}

	_, err := isDelete.Run()
	if err != nil {
		return "", false
	}

	return name, true
}

package cli

import (
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

const (
	Data = "Data"
	Info = "Info"
)

var UpdateItems = []string{Data, Info, ExitKey}

func UpdateSecret(metaData models.MetaData) (string, int, bool) {
	name, ok := chooseSecret(metaData, UpdateSecretKey, updateSecretTemplate)
	if !ok {
		return "", 0, false
	}
	chooseUpdate := promptui.Select{
		Label:     "Data or Info?",
		Items:     UpdateItems,
		Templates: menuTemplate,
		HideHelp:  true,
	}

	i, _, err := chooseUpdate.Run()
	if err != nil {
		return "", 0, false
	}

	if UpdateItems[i] == ExitKey {
		return "", 0, false
	}

	return name, i, true
}

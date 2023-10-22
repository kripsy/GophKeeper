package cli

import (
	"github.com/manifoldco/promptui"
)

const (
	SecretsKey      = "Secrets"
	AddSecretKey    = "Add Secret"
	DeleteSecretKey = "Delete Secret"
	UpdateSecretKey = "Update Secret"
	ExitKey         = "Exit"
)

var MenuTable = []string{SecretsKey, AddSecretKey, UpdateSecretKey, DeleteSecretKey, ExitKey}

func Menu() int {
	defer Clear()
	action := promptui.Select{
		Label:     "Menu",
		Items:     MenuTable,
		Templates: menuTemplate,
		HideHelp:  true,
	}

	id, _, err := action.Run()
	if err != nil {
		return Menu()
	}

	return id
}

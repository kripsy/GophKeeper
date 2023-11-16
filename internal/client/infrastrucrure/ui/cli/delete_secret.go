// Package cli provides command-line interface functionalities for the GophKeeper application.
// It includes methods for interacting with the user for various operations like deleting secrets.
package cli

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

// DeleteSecret allows the user to select and confirm the deletion of a secret.
// It displays a list of secrets to choose from and prompts for confirmation before deletion.
// Returns the name of the secret selected for deletion and a boolean indicating
// whether the deletion process should proceed.
func (c *CLI) DeleteSecret(metaData models.MetaData) (string, bool) {
	// The user is first asked to choose a secret from the list.
	name, ok := c.chooseSecret(metaData, ui.DeleteSecretKey, deleteSecretTemplate)
	if !ok {
		return "", ok
	}

	// The user is then asked to confirm the deletion of the chosen secret.
	isDelete := promptui.Prompt{
		Label:       "Secret will be deleted from the device, are you sure?",
		IsConfirm:   true,
		HideEntered: true,
	}

	// The confirmation prompt is executed and its result is returned.
	_, err := isDelete.Run()
	if err != nil {
		return "", false
	}

	return name, true
}

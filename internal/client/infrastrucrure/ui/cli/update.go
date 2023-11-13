//nolint:gocritic
package cli

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

// UpdateSecret allows the user to select a secret and choose whether to update its data or metadata.
// It first presents a list of secrets to choose from, then offers the choice to update data or metadata.
// Returns the name of the selected secret, the index of the chosen update option,
// and a boolean indicating if a valid selection was made.
func (c *CLI) UpdateSecret(metaData models.MetaData) (string, int, bool) {
	items := append(ui.UpdateItems, ui.ExitKey)

	name, ok := c.chooseSecret(metaData, ui.UpdateSecretKey, updateSecretTemplate)
	if !ok {
		return "", 0, false
	}
	chooseUpdate := promptui.Select{
		Label:        "Data or Info?",
		Items:        items,
		Templates:    menuTemplate,
		HideHelp:     true,
		HideSelected: true,
	}

	i, _, err := chooseUpdate.Run()
	if err != nil {
		return "", 0, false
	}

	if items[i] == ui.ExitKey {
		return "", 0, false
	}

	return name, i, true
}

// Package cli provides command-line interface functionalities for the GophKeeper application.
// It includes methods for displaying the main menu and handling user selections.
package cli

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/manifoldco/promptui"
)

// Menu displays the main application menu and handles the user's selection.
// The menu items are determined based on whether sync storage is available.
// It returns the ID of the selected menu item.
func (c *CLI) Menu(isSyncStorage bool) int {
	defer c.Clear()

	items := ui.MenuTable
	if !isSyncStorage {
		items = ui.LocalMenuTable
	}
	action := promptui.Select{
		Label:        "Menu",
		Items:        items,
		CursorPos:    1,
		Size:         len(items),
		Templates:    menuTemplate,
		HideHelp:     true,
		HideSelected: true,
	}

	id, _, err := action.Run()
	if err != nil {
		return c.Menu(isSyncStorage)
	}

	if !isSyncStorage {
		return id + 1
	}

	return id
}

package cli

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/manifoldco/promptui"
)

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

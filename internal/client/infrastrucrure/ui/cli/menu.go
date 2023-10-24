package cli

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/manifoldco/promptui"
)

func (c CLI) Menu(isLocalStorage bool) int {
	defer c.Clear()

	items := ui.MenuTable
	if isLocalStorage {
		items = ui.LocalMenuTable
	}
	action := promptui.Select{
		Label:     "Menu",
		Items:     items,
		CursorPos: 2,
		Size:      len(items),
		Templates: menuTemplate,
		HideHelp:  true,
	}

	id, _, err := action.Run()
	if err != nil {

		return c.Menu(isLocalStorage)
	}

	if isLocalStorage {
		return id + 1
	}

	return id
}

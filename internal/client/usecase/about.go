package usecase

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// about is a method for displaying information about an application with the version and build date.
func (c *ClientUsecase) about() {
	fmt.Println(promptui.Styler(
		promptui.FGItalic,
		promptui.FGCyan,
	)(c.aboutMsg))
}

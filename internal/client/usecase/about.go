package usecase

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func (c *ClientUsecase) about() {
	defer c.InMenu()
	fmt.Println(promptui.Styler(
		promptui.FGItalic,
		promptui.FGCyan,
	)(c.aboutMsg))
}

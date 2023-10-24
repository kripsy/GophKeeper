package cli

import (
	"errors"
	"github.com/inancgumus/screen"
	"github.com/manifoldco/promptui"
	"os"
)

const (
	dateFormat = "[2006/01/02 15:01]"
)

func (c CLI) checkInterrupt(err error) {
	if errors.Is(err, promptui.ErrInterrupt) {
		c.Exit()
	}
}

func (c CLI) Clear() {
	screen.MoveTopLeft() // не работакет, как и "clean"
	screen.Clear()
}

func (c CLI) Exit() {
	c.Clear()
	os.Exit(1)
}

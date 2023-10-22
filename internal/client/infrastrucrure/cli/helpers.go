package cli

import (
	"github.com/inancgumus/screen"
	"os"
)

const (
	dateFormat = "[2006/01/02 15:01]"
)

func Clear() {
	screen.MoveTopLeft() // не работакет, как и "clean"
	screen.Clear()
}

func Exit() {
	Clear()
	os.Exit(1)
}

package cli

import (
	"errors"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"runtime"
)

const (
	dateFormat = "[2006/01/02 15:01]"
)

func (c CLI) checkInterrupt(err error) {
	if errors.Is(err, promptui.ErrInterrupt) {
		c.Exit()
	}
}

var clearMap map[string]func()

func init() {
	clearMap = make(map[string]func())
	clearMap["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

}

func (c CLI) Clear() {
	value, ok := clearMap[runtime.GOOS] //runtime.GOOS -> linux, windows etc.
	if ok {
		value()
	} else {
		panic("")
	}
}

func (c CLI) Exit() {
	c.Clear()
	os.Exit(1)
}

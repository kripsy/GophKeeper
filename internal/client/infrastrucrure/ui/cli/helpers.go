package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/manifoldco/promptui"
)

const (
	dateFormat = "[2006/01/02 15:01]"
)

func (c *CLI) checkInterrupt(err error) {
	if errors.Is(err, promptui.ErrInterrupt) {
		c.Exit()
	}
}

var clearMapByOS map[string]func() error

func init() {
	clearMapByOS = make(map[string]func() error)
	clearMapByOS["linux"] = clearFunc(exec.Command("clear"))
	clearMapByOS["windows"] = clearFunc(exec.Command("cmd", "/c", "cls"))
	clearMapByOS["default"] = clearFunc(exec.Command("clear"))
}

func (c *CLI) Clear() {
	value, ok := clearMapByOS[runtime.GOOS] //runtime.GOOS -> linux, windows etc.
	var err error
	if ok {
		if err = value(); err != nil {
			c.log.Err(err).Msg("failed get clear func")
		}
		return
	}
	value, ok = clearMapByOS["default"]
	if ok {
		if err = value(); err != nil {
			c.log.Err(err).Msg("failed get clear func")
		}

		return
	}
}

func clearFunc(cmd *exec.Cmd) func() error {
	return func() error {
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error in clear map ", err.Error())
			return err
		}
		return nil
	}
}

func (c *CLI) Exit() {
	c.Clear()
	os.Exit(1)
}

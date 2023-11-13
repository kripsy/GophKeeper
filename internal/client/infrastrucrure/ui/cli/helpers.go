// Package cli provides functionalities for the command-line interface of the GophKeeper application.
// It includes methods for handling terminal interactions such as clearing the screen,
// handling interruptions, and exiting the application.
//
//nolint:gochecknoinits, gochecknoglobals
package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/manifoldco/promptui"
)

// DateFormat is the standard format used for displaying dates in the CLI.
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
	clearMapByOS["windows"] = clearFunc("cmd", "/c", "cls")
	clearMapByOS["default"] = clearFunc("clear")
}

// Clear clears the terminal screen based on the user's operating system.
func (c *CLI) Clear() {
	value, ok := clearMapByOS[runtime.GOOS] // runtime.GOOS -> linux, windows etc.
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

func clearFunc(name string, args ...string) func() error {
	return func() error {
		cmd := exec.Command(name, args...)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("clearFunc: %w", err)
		}

		return nil
	}
}

// Exit clears the terminal screen and exits the application.
func (c *CLI) Exit() {
	c.Clear()
	os.Exit(1)
}

// PrintErr displays an error message in a styled format.
func (c *CLI) PrintErr(err string) {
	fmt.Println(promptui.Styler(
		promptui.FGRed,
		promptui.FGBold,
	)(err))
}

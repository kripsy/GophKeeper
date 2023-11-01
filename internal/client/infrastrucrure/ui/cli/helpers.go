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

func (c CLI) checkInterrupt(err error) {
	if errors.Is(err, promptui.ErrInterrupt) {
		c.Exit()
	}
}

var clearMap map[string]func() error

func init() {
	clearMap = make(map[string]func() error)
	clearMap["linux"] = func() error {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error in clear map ", err.Error())
			return err
		}
		return nil
	}
	clearMap["windows"] = func() error {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error in clear map ", err.Error())
			return err
		}
		return nil
	}
	clearMap["default"] = func() error {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error in clear map ", err.Error())
			return err
		}
		return nil
	}

}

func (c CLI) Clear() {
	value, ok := clearMap[runtime.GOOS] //runtime.GOOS -> linux, windows etc.
	var err error
	if ok {
		err = value()
		if err != nil {
			panic(err.Error())
		}
		return
	}
	value, ok = clearMap["default"]
	if ok {
		err = value()
		if err != nil {
			panic(err.Error())
		}
		return
	} else {
		fmt.Println("I haven't func for clear console...")
		panic("")
	}
}

func (c CLI) Exit() {
	c.Clear()
	os.Exit(1)
}

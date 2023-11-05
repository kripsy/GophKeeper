package cli

import (
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
)

type CLI struct {
	log zerolog.Logger
}

func NewCLI(log zerolog.Logger) *CLI {
	promptui.IconBad = "🌚"
	promptui.IconGood = "🌝"
	return &CLI{
		log: log,
	}
}

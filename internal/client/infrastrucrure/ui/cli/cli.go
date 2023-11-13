// Package cli provides a command-line interface for the GophKeeper application.
// It defines the CLI structure with methods for various user interactions,
// including prompts and data selection.
package cli

import (
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
)

// CLI struct encapsulates the functionality for the command-line interface of the application.
// It includes a logger for logging purposes.
type CLI struct {
	log zerolog.Logger // log is used for logging information and errors.
}

// NewCLI creates a new instance of the CLI with the provided logger.
// It sets custom icons for the prompt UI and returns the CLI instance.
func NewCLI(log zerolog.Logger) *CLI {
	promptui.IconBad = "🌚"  // Custom icon for negative prompts.
	promptui.IconGood = "🌝" // Custom icon for positive prompts.

	return &CLI{
		log: log,
	}
}

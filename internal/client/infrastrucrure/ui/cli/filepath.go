// Package cli provides command-line interface functionalities for the GophKeeper application.
// It includes methods for handling file upload paths and interactively selecting file or directory paths.
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c-bata/go-prompt"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/manifoldco/promptui"
)

const (
	gigabyte = 1_000_000_000.0
	megabyte = 1_000_000.0
	kilobyte = 10_000.0
)

// Constants for predefined directory options.
const (
	CfgDir    = "Into the directory specified in the configuration file"
	CustomDir = "I will specify the directory myself"
)

// UploadFileTo prompts the user to choose a directory for file upload.
// Returns the chosen directory's path and a boolean indicating if a valid selection was made.
func (c *CLI) UploadFileTo(cfgDir string) (string, bool) {
	chooseUpload := promptui.Select{
		Label:        "Where do you want to move the file to?",
		Items:        []string{CfgDir, CustomDir, ui.ExitKey},
		Templates:    menuTemplate,
		HideHelp:     true,
		HideSelected: true,
	}

	_, result, err := chooseUpload.Run()
	if err != nil {
		return "", false
	}
	switch result {
	case CfgDir:
		return cfgDir, true
	case CustomDir:
		var newFilePath string
		c.GetNewFilePath(&newFilePath)

		return newFilePath, true
	case ui.ExitKey:
		return "", false
	}

	return "", false
}

// GetFilePath prompts the user to enter a file path using an interactive auto-complete interface.
// Returns the entered file path.
func (c *CLI) GetFilePath() string {
	defer c.Clear()

	var filePath string
	fmt.Println("Use Tab:")
	prompt.New(
		executor(&filePath),
		completerFile(&filePath),
		prompt.OptionSetExitCheckerOnInput(exit),
		prompt.OptionPrefix("▶ "),
		prompt.OptionSelectedSuggestionBGColor(prompt.Yellow),
		prompt.OptionSelectedDescriptionBGColor(prompt.Blue),

		prompt.OptionSuggestionBGColor(prompt.Blue),
		prompt.OptionScrollbarBGColor(prompt.Blue),
		prompt.OptionDescriptionBGColor(prompt.DarkBlue),

		prompt.OptionPreviewSuggestionTextColor(prompt.Yellow),
	).Run()

	return filePath
}

// GetNewFilePath prompts the user to enter a directory path using an interactive auto-complete interface.
// Modifies the provided string pointer with the entered path.
func (c *CLI) GetNewFilePath(filePath *string) {
	defer c.Clear()
	fmt.Println("Use Tab:")
	prompt.New(
		executor(filePath),
		completerDir(filePath),
		prompt.OptionSetExitCheckerOnInput(exit),
		prompt.OptionPrefix("▶ "),
		prompt.OptionSelectedSuggestionBGColor(prompt.Yellow),
		prompt.OptionSelectedDescriptionBGColor(prompt.Blue),

		prompt.OptionSuggestionBGColor(prompt.Blue),
		prompt.OptionScrollbarBGColor(prompt.Blue),
		prompt.OptionDescriptionBGColor(prompt.DarkBlue),

		prompt.OptionPreviewSuggestionTextColor(prompt.Yellow),
	).Run()
}

func completerFile(path *string) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		var s []prompt.Suggest
		current := d.GetWordBeforeCursor()
		s = append(s, prompt.Suggest{Text: ui.ExitKey, Description: *path})
		s = append(s, prompt.Suggest{Text: "../", Description: "Parent Directory"})

		files, _ := filepath.Glob(current + "*")
		for _, f := range files {
			info, osErr := os.Stat(f)
			if osErr == nil {
				s = append(s, prompt.Suggest{
					Text:        f,
					Description: getDescription(info),
				})
			}
		}

		return prompt.FilterHasPrefix(s, current, true)
	}
}

func completerDir(path *string) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		var s []prompt.Suggest
		current := d.GetWordBeforeCursor()
		s = append(s, prompt.Suggest{Text: ui.ExitKey, Description: *path})
		s = append(s, prompt.Suggest{Text: "../", Description: "Parent Directory"})

		files, _ := filepath.Glob(current + "*")
		for _, f := range files {
			info, osErr := os.Stat(f)
			if osErr == nil {
				if info.IsDir() {
					s = append(s, prompt.Suggest{
						Text:        f,
						Description: getDescription(info),
					})
				}
			}
		}

		return prompt.FilterHasPrefix(s, current, true)
	}
}

func executor(path *string) func(path string) {
	return func(p string) {
		if p != ui.ExitKey {
			*path = p
		}
	}
}

func getDescription(info os.FileInfo) string {
	if info.IsDir() {
		return info.ModTime().Format(dateFormat)
	}

	return fmt.Sprintf("%s %s", convertByte(info.Size()), info.ModTime().Format(dateFormat))
}

func exit(input string, breakLine bool) bool {
	return input == ui.ExitKey && breakLine
}

func convertByte(b int64) string {
	bytes := float64(b)
	switch {
	case bytes >= gigabyte:
		return fmt.Sprintf("%.2fGB", bytes/gigabyte)
	case bytes >= megabyte:
		return fmt.Sprintf("%.2fMB", bytes/megabyte)
	case bytes >= kilobyte:
		return fmt.Sprintf("%.2fKB", bytes/kilobyte)
	default:
		return fmt.Sprintf("%.0fB", bytes)
	}
}

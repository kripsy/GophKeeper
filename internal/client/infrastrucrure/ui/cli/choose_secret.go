// Package cli provides command-line interface functionalities for the GophKeeper application.
// It includes methods for displaying and selecting secrets based on user input.
package cli

import (
	"sort"
	"strings"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

const secretMenuSize = 6

// GetSecret displays a list of secrets and allows the user to select one.
// Returns the selected secret's name and a boolean indicating a valid selection.
func (c *CLI) GetSecret(metaData models.MetaData) (string, bool) {
	return c.chooseSecret(metaData, ui.SecretsKey, chooseSecretTemplate)
}

// chooseSecret displays a selectable list of secrets to the user using the provided template.
// It supports searching within the list.
// Returns the chosen secret's name and a boolean indicating if a selection was made.
func (c *CLI) chooseSecret(metaData models.MetaData, label string, template *promptui.SelectTemplates) (string, bool) {
	dataInfos := getForTemplate(metaData)

	searcher := func(input string, index int) bool {
		di := dataInfos[index]
		name := strings.ReplaceAll(strings.ToLower(di.Name), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:             label,
		Items:             dataInfos,
		Templates:         template,
		StartInSearchMode: true,
		HideHelp:          true,
		HideSelected:      true,
		Size:              secretMenuSize,
		Searcher:          searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return c.GetSecret(metaData)
	}

	_, isSecret := metaData[dataInfos[i].Name]
	if !isSecret {
		return "", false
	}

	return dataInfos[i].Name, true
}

// getForTemplate converts metadata into a slice of TemplateInfo for display.
// The data is sorted by name and includes an exit option.
func getForTemplate(md models.MetaData) []TemplateInfo {
	dataInfo := make([]TemplateInfo, 0, len(md))
	for name, info := range md {
		dataInfo = append(dataInfo, TemplateInfo{
			Name:        name,
			Description: info.Description,
			FileName:    info.FileName,
			DataType:    filemanager.GetTypeName(info.DataType),
			UpdatedAt:   info.UpdatedAt.Format(dateFormat),
		})
	}
	sort.Slice(dataInfo, func(i, j int) bool {
		return dataInfo[i].Name < dataInfo[j].Name
	})

	dataInfo = append(dataInfo, TemplateInfo{Name: ui.ExitKey, Description: "return to Menu", DataType: "◀"})

	return dataInfo
}

// TemplateInfo struct defines the information to be displayed for each secret in the selection menu.
type TemplateInfo struct {
	Name        string  // Name of the secret.
	DataType    string  // Type of the secret (e.g., Note, Login&Password).
	Description string  // Description of the secret.
	FileName    *string // File name associated with the secret.
	UpdatedAt   string  // Last update date of the secret.
}

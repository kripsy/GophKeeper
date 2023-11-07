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

func (c *CLI) GetSecret(metaData models.MetaData) (string, bool) {
	return c.chooseSecret(metaData, ui.SecretsKey, chooseSecretTemplate)
}

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

type TemplateInfo struct {
	Name        string
	DataType    string
	Description string
	FileName    *string
	UpdatedAt   string
}

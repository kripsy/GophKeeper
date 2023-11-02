package cli

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const (
	red    = "red"
	yellow = "yellow"
)

const (
	label = "〘 {{ . | %s }} 〙"

	activeMenu   = "▻「  {{ . | cyan }}  」"
	inactiveMenu = "  [{{ . | blue }}]"

	activeSecret   = "%s  「  {{ .Name | cyan }}  」 {{ .DataType }}  {{ .UpdatedAt | faint }}"
	inactiveSecret = "    [{{ .Name | blue }}] {{ .DataType | faint }}"
	detailsSecret  = `
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{{ "Description:" | faint }}	{{ .Description  }}					
{{ if .FileName }}{{ "FileName:" | faint }}	{{ .FileName }}{{ end }}
`
)

var labelFunc = func(color string) string {
	return fmt.Sprintf(label, color)
}

var activeSecretFunc = func(icon string) string {
	return fmt.Sprintf(activeSecret, icon)
}

var menuTemplate = &promptui.SelectTemplates{
	Label:    labelFunc(yellow),
	Active:   activeMenu,
	Inactive: inactiveMenu,
}

var tryAgainTemplate = &promptui.SelectTemplates{
	Label:    labelFunc(red),
	Active:   activeMenu,
	Inactive: inactiveMenu,
}

var chooseSecretTemplate = &promptui.SelectTemplates{
	Label:    labelFunc(yellow),
	Active:   activeSecretFunc("🔑"),
	Inactive: inactiveSecret,
	Details:  detailsSecret,
}

var deleteSecretTemplate = &promptui.SelectTemplates{
	Label:    labelFunc(red),
	Active:   activeSecretFunc("‼️"),
	Inactive: inactiveSecret,
	Details:  detailsSecret,
}

var updateSecretTemplate = &promptui.SelectTemplates{
	Label:    labelFunc(red),
	Active:   activeSecretFunc("♻️"),
	Inactive: inactiveSecret,
	Details:  detailsSecret,
}

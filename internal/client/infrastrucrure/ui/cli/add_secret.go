package cli

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/validation"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

func (c *CLI) ChooseSecretType() (int, bool) {
	defer c.Clear()
	items := filemanager.DataTypeTable
	items = append(items, ui.ExitKey)
	action := promptui.Select{
		Label:        ui.AddSecretKey,
		Items:        items,
		Templates:    menuTemplate,
		HideHelp:     true,
		HideSelected: true,
	}

	id, _, err := action.Run()
	if err != nil {
		c.checkInterrupt(err)

		return 0, false
	}
	if items[id] == ui.ExitKey {
		return 0, false
	}

	return id, true
}

func (c *CLI) AddNote() (filemanager.Note, error) {
	note := promptui.Prompt{
		Label:       "Note",
		HideEntered: true,
	}

	text, err := note.Run()
	if err != nil {
		c.checkInterrupt(err)

		return filemanager.Note{}, fmt.Errorf("%w", err)
	}

	return filemanager.Note{
		Text: text,
	}, nil
}

func (c *CLI) AddBasicAuth() (filemanager.BasicAuth, error) {
	log := promptui.Prompt{
		Label:       "Login",
		HideEntered: true,
	}

	pass := promptui.Prompt{
		Label:       "Password",
		HideEntered: true,
	}

	login, err := log.Run()
	if err != nil {
		c.checkInterrupt(err)

		return filemanager.BasicAuth{}, fmt.Errorf("%w", err)
	}
	password, err := pass.Run()
	if err != nil {
		return filemanager.BasicAuth{}, fmt.Errorf("%w", err)
	}

	return filemanager.BasicAuth{
		Login:    login,
		Password: password,
	}, nil
}

func (c *CLI) AddCard() (filemanager.CardData, error) {
	cardNum := promptui.Prompt{
		Label:       "Card Number",
		Validate:    validation.ValidateCardNumber,
		HideEntered: true,
	}

	cardDate := promptui.Prompt{
		Label:       "Card Date",
		HideEntered: true,
	}

	cardCvv := promptui.Prompt{
		Label:       "CVV",
		HideEntered: true,
		Validate:    validation.ValidateCVV,
		Mask:        '⏀',
	}

	number, err := cardNum.Run()
	if err != nil {
		return filemanager.CardData{}, fmt.Errorf("%w", err)
	}

	date, err := cardDate.Run()
	if err != nil {
		return filemanager.CardData{}, fmt.Errorf("%w", err)
	}

	cvv, err := cardCvv.Run()
	if err != nil {
		return filemanager.CardData{}, fmt.Errorf("%w", err)
	}

	return filemanager.CardData{
		Number: number,
		Date:   date,
		CVV:    cvv,
	}, nil
}

func (c *CLI) AddMetaInfo() (models.DataInfo, error) {
	dataName := promptui.Prompt{
		Label:       "Secret name",
		HideEntered: true,
	}

	dataDescription := promptui.Prompt{
		Label:       "Secret Description",
		HideEntered: true,
	}

	name, err := dataName.Run()
	if err != nil {
		return models.DataInfo{}, fmt.Errorf("%w", err)
	}

	description, err := dataDescription.Run()
	if err != nil {
		return models.DataInfo{}, fmt.Errorf("%w", err)
	}

	return models.DataInfo{
		Name:        name,
		Description: description,
	}, nil
}

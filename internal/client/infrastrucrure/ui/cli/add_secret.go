package cli

import (
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

func (c CLI) ChooseSecretType() (int, bool) {
	defer c.Clear()
	items := filemanager.DataTypeTable
	items = append(items, ui.ExitKey)
	action := promptui.Select{
		Label:     ui.AddSecretKey,
		Items:     items,
		Templates: menuTemplate,
		HideHelp:  true,
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

func (c CLI) AddNote() (filemanager.Note, error) {
	note := promptui.Prompt{
		Label: "Note",
	}

	text, err := note.Run()
	if err != nil {
		c.checkInterrupt(err)

		return filemanager.Note{}, err
	}

	return filemanager.Note{
		Text: text,
	}, nil
}

func (c CLI) AddBasicAuth() (filemanager.BasicAuth, error) {
	log := promptui.Prompt{
		Label: "Login",
	}

	pass := promptui.Prompt{
		Label: "Password",
	}

	login, err := log.Run()
	if err != nil {
		c.checkInterrupt(err)

		return filemanager.BasicAuth{}, err
	}
	password, err := pass.Run()
	if err != nil {
		return filemanager.BasicAuth{}, err
	}

	return filemanager.BasicAuth{
		Login:    login,
		Password: password,
	}, nil
}

func (c CLI) AddCard() (filemanager.CardData, error) {
	cardNum := promptui.Prompt{
		Label:    "Card Number",
		Validate: validateCardNumber,
	}

	cardDate := promptui.Prompt{
		Label: "Card Date",
	}

	cardCvv := promptui.Prompt{
		Label:       "CVV",
		HideEntered: true,
		Validate:    validateCVV,
		Mask:        '⏀',
	}

	number, err := cardNum.Run()
	if err != nil {
		return filemanager.CardData{}, err
	}

	date, err := cardDate.Run()
	if err != nil {
		return filemanager.CardData{}, err
	}

	cvv, err := cardCvv.Run()
	if err != nil {
		return filemanager.CardData{}, err
	}

	return filemanager.CardData{
		Number: number,
		Date:   date,
		CVV:    cvv,
	}, nil
}

func (c CLI) AddMetaInfo() (models.DataInfo, error) {
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
		return models.DataInfo{}, err
	}

	description, err := dataDescription.Run()
	if err != nil {
		return models.DataInfo{}, err
	}

	return models.DataInfo{
		Name:        name,
		Description: description,
	}, nil
}

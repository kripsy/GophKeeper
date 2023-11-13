// Package cli provides command-line interface functionalities for the GophKeeper application.
// It includes methods for interacting with the user, such as choosing secret types, adding notes,
// basic authentication data, card data, and metadata.
package cli

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/validation"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

// ChooseSecretType prompts the user to choose a type of secret from a predefined list.
// Returns the chosen secret type's ID and a boolean indicating if a valid selection was made.
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

// AddNote prompts the user to enter a note.
// Returns a Note object containing the entered text or an error if the operation fails.
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

// AddBasicAuth prompts the user to enter login and password credentials.
// Returns a BasicAuth object containing the entered credentials or an error if the operation fails.
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

// AddCard prompts the user to enter card details (number, date, and CVV).
// Returns a CardData object containing the entered details or an error if the operation fails.
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

// AddMetaInfo prompts the user to enter metadata for a secret (name and description).
// Returns a DataInfo object containing the entered metadata or an error if the operation fails.
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

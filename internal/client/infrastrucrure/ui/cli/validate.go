//nolint:gochecknoglobals
package cli

import (
	"errors"
	"strconv"

	"github.com/manifoldco/promptui"
)

const (
	userSize     = 6
	passwordSize = 6
	cvvSize      = 3
	CardNumSize  = 16
)

var errStyle = promptui.Styler(promptui.BGRed, promptui.FGBold, promptui.FGBlack)

func validatePassword(input string) error {
	if len(input) < userSize {
		return errors.New(errStyle("Password must have more than 6 characters"))
	}

	return nil
}

func validateUsername(input string) error {
	if len(input) < passwordSize {
		return errors.New(errStyle("Username must have more than 6 characters"))
	}

	return nil
}

func validateCVV(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return errors.New("CVV must  be a number")
	}

	if len(input) != cvvSize {
		return errors.New("CVV must contain 3 digits")
	}

	return nil
}

func validateCardNumber(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return errors.New("Card number must  be a number ツ")
	}

	if len(input) > CardNumSize {
		return errors.New("Card number cannot contain more than 16 digits")
	}

	return nil
}

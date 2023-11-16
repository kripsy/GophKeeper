//nolint:gochecknoglobals
package validation

import (
	"errors"

	"github.com/manifoldco/promptui"
)

var errStyle = promptui.Styler(promptui.BGRed, promptui.FGBold, promptui.FGBlack)

var (
	errValidatePassword         = errors.New(errStyle("password must have more than 6 characters"))
	errValidateUsername         = errors.New(errStyle("username must have more than 6 characters"))
	errValidateCVVMustBeNum     = errors.New(errStyle("cVV must  be a number"))
	errValidateCVVSize          = errors.New(errStyle("cVV must contain 3 digits"))
	errValidateCardNumberNotNum = errors.New(errStyle("card number must  be a number ツ"))
	errValidateCardNumSize      = errors.New(errStyle("card number cannot contain more than 16 digits"))
)

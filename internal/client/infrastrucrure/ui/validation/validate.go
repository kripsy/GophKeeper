package validation

import (
	"fmt"
	"strconv"
)

const (
	userSize     = 6
	passwordSize = 6
	cvvSize      = 3
	CardNumSize  = 16
)

func ValidatePassword(input string) error {
	if len(input) < userSize {
		return fmt.Errorf("%w", errValidatePassword)
	}

	return nil
}

func ValidateUsername(input string) error {
	if len(input) < passwordSize {
		return fmt.Errorf("%w", errValidateUsername)
	}

	return nil
}

func ValidateCVV(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return fmt.Errorf("%w", errValidateCVVMustBeNum)
	}

	if len(input) != cvvSize {
		return fmt.Errorf("%w", errValidateCVVSize)
	}

	return nil
}

func ValidateCardNumber(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return fmt.Errorf("%w", errValidateCardNumberNotNum)
	}

	if len(input) > CardNumSize {
		return fmt.Errorf("%w", errValidateCardNumSize)
	}

	return nil
}

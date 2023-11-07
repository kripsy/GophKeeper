package cli

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

func validatePassword(input string) error {
	if len(input) < userSize {
		return fmt.Errorf("%w", errValidatePassword)
	}

	return nil
}

func validateUsername(input string) error {
	if len(input) < passwordSize {
		return fmt.Errorf("%w", errValidateUsername)
	}

	return nil
}

func validateCVV(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return fmt.Errorf("%w", errValidateCVVMustBeNum)
	}

	if len(input) != cvvSize {
		return fmt.Errorf("%w", errValidateCVVSize)
	}

	return nil
}

func validateCardNumber(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return fmt.Errorf("%w", errValidateCardNumberNotNum)
	}

	if len(input) > CardNumSize {
		return fmt.Errorf("%w", errValidateCardNumSize)
	}

	return nil
}

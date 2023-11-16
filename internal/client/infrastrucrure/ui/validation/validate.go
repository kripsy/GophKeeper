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

// ValidatePassword checks if the provided username meets the required length.
func ValidatePassword(input string) error {
	if len(input) < userSize {
		return fmt.Errorf("%w", errValidatePassword)
	}

	return nil
}

// ValidateUsername checks if the provided username meets the required length.
func ValidateUsername(input string) error {
	if len(input) < passwordSize {
		return fmt.Errorf("%w", errValidateUsername)
	}

	return nil
}

// ValidateCVV checks if the provided CVV is a numeric string with the required length.
func ValidateCVV(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return fmt.Errorf("%w", errValidateCVVMustBeNum)
	}

	if len(input) != cvvSize {
		return fmt.Errorf("%w", errValidateCVVSize)
	}

	return nil
}

// ValidateCardNumber checks if the provided credit card number is a numeric string with a valid length.
func ValidateCardNumber(input string) error {
	if _, err := strconv.Atoi(input); err != nil {
		return fmt.Errorf("%w", errValidateCardNumberNotNum)
	}

	if len(input) > CardNumSize {
		return fmt.Errorf("%w", errValidateCardNumSize)
	}

	return nil
}

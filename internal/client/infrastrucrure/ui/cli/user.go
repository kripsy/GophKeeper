package cli

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/validation"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

// GetUser prompts the user to enter their username and password.
// It validates the input and returns a User object or an error if the input is invalid.
func (c *CLI) GetUser() (models.User, error) {
	username := promptui.Prompt{
		Label:       "Username",
		Validate:    validation.ValidateUsername,
		HideEntered: true,
	}
	pass := promptui.Prompt{
		Label:       "Password",
		Validate:    validation.ValidatePassword,
		HideEntered: true,
		Mask:        '#',
	}

	user, err := username.Run()
	if err != nil {
		c.checkInterrupt(err)

		return models.User{}, fmt.Errorf("%w", err)
	}
	password, err := pass.Run()
	if err != nil {
		c.checkInterrupt(err)

		return models.User{}, fmt.Errorf("%w", err)
	}

	return models.User{
		Username: user,
		Password: password,
	}, nil
}

// GetRepeatedPassword prompts the user to repeat their password for confirmation.
// It validates the input and returns the confirmed password or an error if the input is invalid.
func (c *CLI) GetRepeatedPassword() (string, error) {
	pass := promptui.Prompt{
		Label:       "Repeat Password",
		Validate:    validation.ValidatePassword,
		HideEntered: true,
		Mask:        '#',
	}

	password, err := pass.Run()
	if err != nil {
		c.checkInterrupt(err)

		return "", fmt.Errorf("%w", err)
	}

	return password, nil
}

// TryAgain offers the user options to try logging in again, register, or exit the application.
// It handles the user's choice and returns a boolean indicating whether to retry the login.
func (c *CLI) TryAgain() bool {
	defer c.Clear()
	action := promptui.Select{
		Label:        "This user does not exist",
		Items:        []string{"Try Again", "Register", ui.ExitKey},
		Templates:    tryAgainTemplate,
		HideHelp:     true,
		HideSelected: true,
	}
	_, result, err := action.Run()
	if err != nil {
		return c.TryAgain()
	}

	switch result {
	case ui.ExitKey:
		c.Exit()
	case "Try Again":
		return true
	}

	return false
}

// IsSyncStorage prompts the user to choose whether to synchronize their secrets across devices.
// It returns a boolean indicating the user's preference for syncing.
func (c *CLI) IsSyncStorage() bool {
	isLocal := promptui.Prompt{
		Label:       "Do you want to synchronize your secrets across devices?",
		HideEntered: true,
		IsConfirm:   true,
	}

	_, err := isLocal.Run()

	return err == nil
}

package cli

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/validation"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

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

func (c *CLI) IsSyncStorage() bool {
	isLocal := promptui.Prompt{
		Label:       "Do you want to synchronize your secrets across devices?",
		HideEntered: true,
		IsConfirm:   true,
	}

	_, err := isLocal.Run()

	return err == nil
}

package cli

import (
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/manifoldco/promptui"
)

func GetUser() (models.User, error) {
	username := promptui.Prompt{
		Label:       "Username",
		Validate:    validateUsername,
		HideEntered: true,
	}

	pass := promptui.Prompt{
		Label:       "Password",
		Validate:    validatePassword,
		HideEntered: true,
		Mask:        '⏀',
	}

	user, err := username.Run()
	if err != nil {
		return models.User{}, err
	}
	password, err := pass.Run()
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Username: user,
		Password: password,
	}, nil
}

func TryAgain() bool {
	defer Clear()
	action := promptui.Select{
		Label:     "This user does not exist",
		Items:     []string{"Try Again", "Register", ExitKey},
		Templates: tryAgainTemplate,
		HideHelp:  true,
	}
	_, result, err := action.Run()
	if err != nil {
		return TryAgain()
	}

	switch result {
	case ExitKey:
		Exit()
	case "Try Again":
		return true
	}

	return false
}

func IsLocalStorage() bool {
	isLocal := promptui.Prompt{
		Label:     "Do you want to synchronize your secrets across devices?",
		IsConfirm: true,
	}

	_, err := isLocal.Run()
	if err != nil {
		return true
	}

	return false
}

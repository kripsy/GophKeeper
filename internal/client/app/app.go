package app

import (
	"fmt"
	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/cli"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/models"
)

type Application struct {
	dataPath    string
	uploadPath  string
	userData    *models.UserData
	fileManager *filemanager.FileManager
}

func NewApplication(cfg config.Config) (*Application, error) {

	return &Application{
		dataPath:   cfg.StoragePath,
		uploadPath: cfg.UploadPath,
		userData:   &models.UserData{},
	}, nil
}

func (a *Application) PrepareApp() {
	if err := a.setUser(); err != nil {
		panic(err) //todo errs
	}

	if err := a.setFileManager(); err != nil {
		panic(err)
	}

}

func (a *Application) Run() {
	defer cli.Clear()
	a.inMenu()
}

func (a *Application) setUser() error {
	var err error
	userAuth, err := filemanager.NewUserAuth(a.dataPath)
	if err != nil {
		return err
	}

	for {
		a.userData.User, err = cli.GetUser()
		if err != nil {
			fmt.Println(err)
			// todo проверка на конкретную ошкибку
			continue
			//	return err
		}
		// todo повторный ввод пароля при регистрации
		if userAuth.IsUserNotExisting(a.userData.User.GetDir(a.dataPath)) {
			if cli.TryAgain() {
				continue
			}

			a.userData.Meta, err = userAuth.CreateUser(&a.userData.User, cli.IsLocalStorage())
			if err != nil {
				return err
			}
		} else {
			a.userData.Meta, err = userAuth.GetUser(&a.userData.User)
			if err != nil {
				fmt.Println(err)
				// todo проверка на конкретную ошкибку
				continue
				//	return err
			}
		}

		return nil
	}
}

func (a *Application) setFileManager() error {
	fileManager, err := filemanager.NewFileManager(
		a.dataPath,
		a.uploadPath,
		a.userData.User.GetDir(a.dataPath),
		a.userData.Meta,
		a.userData.User.Key,
	)
	if err != nil {
		return err
	}
	a.fileManager = fileManager

	return nil
}

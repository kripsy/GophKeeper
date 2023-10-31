package usecase

import (
	"context"
	"fmt"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/client/grpc"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
	"os"
)

type ClientUsecase struct {
	dataPath      string
	uploadPath    string
	aboutMsg      string
	serverAddress string
	userData      *models.UserData
	grpc          grpc.Client
	fileManager   filemanager.FileStorage
	ui            ui.UserInterface
	log           zerolog.Logger
}

func NewUsecase(
	dataPath string,
	uploadPath string,
	aboutMsg string,
	serverAddress string,
	ui ui.UserInterface,
	log zerolog.Logger,
) *ClientUsecase {
	return &ClientUsecase{
		dataPath:   dataPath,
		uploadPath: uploadPath,
		aboutMsg:   aboutMsg,

		userData: &models.UserData{},
		grpc:     grpc.NewClient(serverAddress, log),
		ui:       ui,
		log:      log,
	}
}

func (c *ClientUsecase) SetUser() error {
	var err error
	userAuth, err := filemanager.NewUserAuth(c.dataPath)
	if err != nil {
		return err
	}

	for {
		c.userData.User, err = c.ui.GetUser()
		if err != nil {
			// todo проверка на конкретную ошкибку
			fmt.Println(err)
			continue
		}

		// todo повторный ввод пароля при регистрации
		if userAuth.IsUserNotExisting(c.userData.User.GetDir(c.dataPath)) {
			if err := c.handleUserRegistration(userAuth); err != nil {
				return err
			}
			return nil
		} else {
			if err := c.handleUserLogin(userAuth); err != nil {
				fmt.Println(err)
				continue
			}
			return nil
		}
	}
}

func (c *ClientUsecase) handleUserRegistration(userAuth *filemanager.UserAuth) error {
	if c.ui.TryAgain() {
		return nil
	}

	if c.grpc.TryToConnect() {
		fmt.Println("Could not connect to the server, only local registration is available")
	}

	repeatedPass, err := c.ui.GetRepeatedPassword()
	if err != nil {

	}

	if c.userData.User.Password != repeatedPass {
		fmt.Println("Password mismatch")
		os.Exit(9)
		return nil
	}

	var isLocalStorage bool
	if c.grpc.IsAvailable() {
		isLocalStorage = c.ui.IsLocalStorage()
	}
	if !isLocalStorage {
		hash, err := c.userData.User.GetHashedPass()
		if err != nil {
			c.log.Err(err).Msg("failed get hashed password")
			return err
		}

		resp, err := c.grpc.Register(context.Background(), &pb.AuthRequest{
			Username: c.userData.User.Username,
			Password: hash,
		})
		if err != nil {
			c.log.Err(err).Str("user", c.userData.User.Username).Msg("failed register user")
			return err
		}
		c.userData.User.Token = resp.Token
	}

	meta, err := userAuth.CreateUser(&c.userData.User, isLocalStorage)
	if err != nil {
		return err
	}
	c.userData.Meta = meta

	return nil
}

func (c *ClientUsecase) handleUserLogin(userAuth *filemanager.UserAuth) error {
	var err error
	switch c.grpc.TryToConnect() {
	case true:
		hash, err := c.userData.User.GetHashedPass()
		if err != nil {
			c.log.Err(err).Msg("failed get hashed password")
			return err
		}

		resp, err := c.grpc.Login(context.Background(), &pb.AuthRequest{
			Username: c.userData.User.Username,
			Password: hash,
		})
		if err != nil {
			c.log.Err(err).Str("user", c.userData.User.Username).Msg("failed login user")
		}
		if resp != nil {
			c.userData.User.Token = resp.Token
		}
	case false:
		fmt.Println("Could not connect to the server, data synchronization will not be available")
	}

	if err == nil {
		c.userData.Meta, err = userAuth.GetUser(&c.userData.User)
	}

	return err
}

func (c *ClientUsecase) SetFileManager() error {
	fileManager, err := filemanager.NewFileManager(
		c.dataPath,
		c.uploadPath,
		c.userData.User.GetDir(c.dataPath),
		c.userData.Meta,
		c.userData.User.Key,
	)
	if err != nil {
		return err
	}
	c.fileManager = fileManager

	return nil
}

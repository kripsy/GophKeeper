package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/client/grpc"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
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
			c.log.Err(err).Msg("failed get user")
			continue
		}

		if userAuth.IsUserNotExisting(c.userData.User.GetDir(c.dataPath)) {
			if c.grpc.TryToConnect() && c.checkUserOnServer(userAuth) {
				return nil
			}
			if c.ui.TryAgain() {
				continue
			}
			if err = c.handleUserRegistration(userAuth); err != nil {
				return err
			}
			return nil
		} else {
			if err = c.handleUserLogin(userAuth); err != nil {
				c.log.Err(err).Msg("")
				continue
			}
			return nil
		}
	}
}

func (c *ClientUsecase) checkUserOnServer(userAuth *filemanager.UserAuth) bool {
	hash, err := c.userData.User.GetHashedPass()
	if err != nil {
		c.log.Info().Err(err).Msg("failed get hashed password")
		return false
	}

	err = c.grpc.Login(c.userData.User.Username, hash)
	if err != nil {
		c.log.Info().Str("user", c.userData.User.Username).Msg("failed login user")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	syncKey := uuid.New().String()

	if err = c.blockSync(ctx, syncKey); err != nil {
		return false
	}

	serverMeta, err := c.downloadServerMeta(ctx, syncKey)
	if err != nil {
		c.log.Info().Err(err).Msg("failed download server meta data")
		return false
	}
	if serverMeta.Username != c.userData.User.Username {
		c.log.Info().Msg("failed to verify user")

		return false
	}

	meta, err := userAuth.CreateUser(&c.userData.User, false)
	if err != nil {
		return false
	}
	c.userData.Meta = meta
	if err = c.grpc.ApplyChanges(ctx, syncKey); err != nil {
		c.log.Info().Err(err).Msg("failed apply changes")

		return false
	}

	return true
}

func (c *ClientUsecase) handleUserRegistration(userAuth *filemanager.UserAuth) error {
	if c.grpc.IsNotAvailable() {
		fmt.Println("Could not connect to the server, only local registration is available")
	}

	repeatedPass, err := c.ui.GetRepeatedPassword()
	if err != nil {
		c.log.Err(err).Msg("failed repeated pass")
		return err
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

		err = c.grpc.Register(c.userData.User.Username, hash)
		if err != nil {
			c.log.Err(err).Str("user", c.userData.User.Username).Msg("failed register user")
			return err
		}
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

		err = c.grpc.Login(c.userData.User.Username, hash)
		if err != nil {
			c.log.Err(err).Str("user", c.userData.User.Username).Msg("failed login user")
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

package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/client/grpc"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
)

// Define a custom error for password mismatch.
var errPasswordMismatch = errors.New("password mismatch")

// ClientUsecase structure encapsulates the business logic for the GophKeeper client.
type ClientUsecase struct {
	dataPath    string
	uploadPath  string
	aboutMsg    string
	userData    *models.UserData
	grpc        grpc.Client
	fileManager filemanager.FileStorage
	ui          ui.UserInterface
	log         zerolog.Logger
}

// NewUsecase creates a new instance of the GophKeeper client use case.
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
		userData:   &models.UserData{},
		grpc:       grpc.NewClient(serverAddress, log),
		ui:         ui,
		log:        log,
	}
}

// SetUser handles the user registration or login process.
func (c *ClientUsecase) SetUser() error {
	var err error
	userAuth, err := filemanager.NewUserAuth(c.dataPath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for {
		c.userData.User, err = c.ui.GetUser()
		if err != nil {
			c.log.Err(err).Msg("failed get user")

			continue
		}

		if userAuth.IsUseExisting(c.userData.User.GetDir(c.dataPath)) {
			if err = c.handleUserLogin(userAuth); err != nil {
				c.log.Err(err).Msg("failed handle user login")

				continue
			}

			return nil
		}
		if c.grpc.TryToConnect() && c.checkUserOnServer(userAuth) {
			return nil
		}
		if c.ui.TryAgain() {
			continue
		}

		return c.handleUserRegistration(userAuth)
	}
}

// checkUserOnServer checks if the user exists on the server and handles the registration process if needed.
func (c *ClientUsecase) checkUserOnServer(userAuth filemanager.Auth) bool {
	hash, err := c.userData.User.GetHashedPass()
	if err != nil {
		c.log.Info().Err(err).Msg("failed get hashed password")

		return false
	}

	err = c.grpc.Login(c.userData.User.Username, hash)
	if err != nil {
		c.log.Info().Str("user", c.userData.User.Username).Msg("failed login user")

		return false
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

	meta, err := userAuth.CreateUser(&c.userData.User, true)
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

// handleUserRegistration handles the user registration process.
func (c *ClientUsecase) handleUserRegistration(userAuth filemanager.Auth) error {
	if c.grpc.IsNotAvailable() {
		c.ui.PrintErr("Could not connect to the server, only local registration is available")
	}

	repeatedPass, err := c.ui.GetRepeatedPassword()
	if err != nil {
		c.log.Err(err).Msg("failed repeated pass")

		return fmt.Errorf("%w", err)
	}

	if c.userData.User.Password != repeatedPass {
		c.ui.PrintErr(errPasswordMismatch.Error())

		return fmt.Errorf("%w", errPasswordMismatch)
	}

	var isSyncStorage bool
	if c.grpc.IsAvailable() {
		isSyncStorage = c.ui.IsSyncStorage()
	}
	if isSyncStorage {
		hash, err := c.userData.User.GetHashedPass()
		if err != nil {
			c.log.Err(err).Msg("failed get hashed password")

			return fmt.Errorf("%w", err)
		}

		err = c.grpc.Register(c.userData.User.Username, hash)
		if err != nil {
			c.log.Err(err).Str("user", c.userData.User.Username).Msg("failed register user")

			return fmt.Errorf("%w", err)
		}
	}

	meta, err := userAuth.CreateUser(&c.userData.User, isSyncStorage)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	c.userData.Meta = meta

	return nil
}

// handleUserLogin handles the user login process.
func (c *ClientUsecase) handleUserLogin(userAuth filemanager.Auth) error {
	var err error
	switch c.grpc.TryToConnect() {
	case true:
		var hash string
		hash, err = c.userData.User.GetHashedPass()
		if err != nil {
			c.log.Err(err).Msg("failed get hashed password")

			return fmt.Errorf("%w", err)
		}

		err = c.grpc.Login(c.userData.User.Username, hash)
		if err != nil {
			c.log.Err(err).Str("user", c.userData.User.Username).Msg("failed login user")
		}

	case false:
		fmt.Println("Could not connect to the server, data synchronization will not be available")
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if c.userData.Meta, err = userAuth.GetUser(&c.userData.User); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// SetFileManager sets up the users file manager for the client use case.
func (c *ClientUsecase) SetFileManager() error {
	fileManager, err := filemanager.NewFileManager(
		c.dataPath,
		c.uploadPath,
		c.userData.User.GetDir(c.dataPath),
		c.userData.Meta,
		c.userData.User.Key,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	c.fileManager = fileManager

	return nil
}

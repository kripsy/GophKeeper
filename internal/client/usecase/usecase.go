package usecase

import (
	"context"
	"fmt"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
)

type ClientUsecase struct {
	dataPath    string
	uploadPath  string
	aboutMsg    string
	userData    *models.UserData
	grpc        pb.GophKeeperServiceClient
	fileManager *filemanager.FileManager
	ui          ui.UserInterface
	log         zerolog.Logger
}

func NewUsecase(
	dataPath string,
	uploadPath string,
	aboutMsg string,
	grpc pb.GophKeeperServiceClient,
	ui ui.UserInterface,
	log zerolog.Logger,
) *ClientUsecase {
	return &ClientUsecase{
		dataPath:   dataPath,
		uploadPath: uploadPath,
		aboutMsg:   aboutMsg,
		userData:   &models.UserData{},
		grpc:       grpc,
		ui:         ui,
		log:        log,
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
			fmt.Println(err)
			// todo проверка на конкретную ошкибку
			continue
			//	return err
		}
		// todo повторный ввод пароля при регистрации
		if userAuth.IsUserNotExisting(c.userData.User.GetDir(c.dataPath)) {
			if c.ui.TryAgain() {
				continue
			}

			isLocalStorage := c.ui.IsLocalStorage()
			if !isLocalStorage {
				hash, err := c.userData.User.GetHashedPass()
				if err != nil {
					c.log.Err(err).Msg("failed get hashed password")
					continue
				}

				resp, err := c.grpc.Register(context.Background(), &pb.AuthRequest{
					Username: c.userData.User.Username,
					Password: hash,
				})
				if err != nil {
					c.log.Err(err).Str("user", c.userData.User.Username).Msg("failed register user")

					continue
				}
				c.userData.User.Token = resp.Token
			}

			c.userData.Meta, err = userAuth.CreateUser(&c.userData.User, isLocalStorage)
			if err != nil {
				return err
			}
		} else {

			hash, err := c.userData.User.GetHashedPass()
			if err != nil {
				c.log.Err(err).Msg("failed get hashed password")
				continue
			}

			resp, err := c.grpc.Login(context.Background(), &pb.AuthRequest{
				Username: c.userData.User.Username,
				Password: hash,
			})
			if err != nil {
				fmt.Println("unsuccessful attempt to connect to the server, data synchronization will not be available")
				c.log.Err(err).Str("user", c.userData.User.Username).Msg("failed login user")
			}
			if resp != nil {
				c.userData.User.Token = resp.Token
			}
			c.userData.Meta, err = userAuth.GetUser(&c.userData.User)
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

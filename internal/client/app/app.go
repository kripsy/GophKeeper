package app

import (
	"context"
	"crypto/tls"
	"fmt"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/cli"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Application struct {
	dataPath    string
	uploadPath  string
	userData    *models.UserData
	fileManager *filemanager.FileManager
	grpc        pb.GophKeeperServiceClient
	ui          ui.UserInterface
	buildInfo   BuildInfo
	log         zerolog.Logger
}

type BuildInfo struct {
	BuildVersion string
	BuildDate    string
}

func NewApplication(cfg config.Config, bi BuildInfo, log zerolog.Logger) (*Application, error) {
	conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		log.Info().Err(err).Msg("failed connect to server")

	}
	//defer conn.Close()

	client := pb.NewGophKeeperServiceClient(conn)

	return &Application{
		dataPath:   cfg.StoragePath,
		uploadPath: cfg.UploadPath,
		userData:   &models.UserData{},
		grpc:       client,
		ui:         cli.NewCLI(log),
		buildInfo:  bi,
		log:        log,
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
	defer a.ui.Clear()
	a.inMenu()
}

func (a *Application) setUser() error {
	var err error
	userAuth, err := filemanager.NewUserAuth(a.dataPath)
	if err != nil {
		return err
	}

	for {
		a.userData.User, err = a.ui.GetUser()
		if err != nil {
			fmt.Println(err)
			// todo проверка на конкретную ошкибку
			continue
			//	return err
		}
		// todo повторный ввод пароля при регистрации
		if userAuth.IsUserNotExisting(a.userData.User.GetDir(a.dataPath)) {
			if a.ui.TryAgain() {
				continue
			}

			isLocalStorage := a.ui.IsLocalStorage()
			if !isLocalStorage {
				hash, err := a.userData.User.GetHashedPass()
				if err != nil {
					a.log.Err(err).Msg("failed get hashed password")
					continue
				}

				resp, err := a.grpc.Register(context.Background(), &pb.AuthRequest{
					Username: a.userData.User.Username,
					Password: hash,
				})
				if err != nil {
					a.log.Err(err).Str("user", a.userData.User.Username).Msg("failed register user")

					continue
				}
				a.userData.User.Token = resp.Token
			}

			a.userData.Meta, err = userAuth.CreateUser(&a.userData.User, isLocalStorage)
			if err != nil {
				return err
			}
		} else {

			hash, err := a.userData.User.GetHashedPass()
			if err != nil {
				a.log.Err(err).Msg("failed get hashed password")
				continue
			}

			resp, err := a.grpc.Login(context.Background(), &pb.AuthRequest{
				Username: a.userData.User.Username,
				Password: hash,
			})
			if err != nil {
				fmt.Println("unsuccessful attempt to connect to the server, data synchronization will not be available")
				a.log.Err(err).Str("user", a.userData.User.Username).Msg("failed login user")
			}
			if resp != nil {
				a.userData.User.Token = resp.Token
			}
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

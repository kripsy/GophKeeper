package app

import (
	"crypto/tls"
	"github.com/kripsy/GophKeeper/internal/client/usecase"

	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/cli"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Application struct {
	usecase *usecase.ClientUsecase
	log     zerolog.Logger
}

func NewApplication(cfg config.Config, bi BuildInfo, log zerolog.Logger) (*Application, error) {
	conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		log.Info().Err(err).Msg("failed connect to server")

	}
	//defer conn.Close()

	client := pb.NewGophKeeperServiceClient(conn)

	return &Application{
		usecase: usecase.NewUsecase(
			cfg.StoragePath,
			cfg.UploadPath,
			about(bi),
			client,
			cli.NewCLI(log),
			log,
		),
		log: log,
	}, nil
}

func (a *Application) PrepareApp() {
	if err := a.usecase.SetUser(); err != nil {
		panic(err) //todo errs
	}

	if err := a.usecase.SetFileManager(); err != nil {
		panic(err)
	}

}

func (a *Application) Run() {
	a.usecase.InMenu()
}

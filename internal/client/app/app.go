package app

import (
	"github.com/kripsy/GophKeeper/internal/client/usecase"

	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/cli"
	"github.com/rs/zerolog"
)

type Application struct {
	usecase *usecase.ClientUsecase
	log     zerolog.Logger
}

func NewApplication(cfg config.Config, bi BuildInfo, log zerolog.Logger) (*Application, error) {

	return &Application{
		usecase: usecase.NewUsecase(
			cfg.StoragePath,
			cfg.UploadPath,
			about(bi),
			cfg.ServerAddress,
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

package app

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/cli"
	"github.com/kripsy/GophKeeper/internal/client/usecase"
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

func (a *Application) PrepareApp() error {
	if err := a.usecase.SetUser(); err != nil {
		return fmt.Errorf("SetUser: %w", err)
	}

	if err := a.usecase.SetFileManager(); err != nil {
		return fmt.Errorf("SetFileManager: %w", err)
	}

	return nil
}

func (a *Application) Run() {
	a.usecase.InMenu()
}

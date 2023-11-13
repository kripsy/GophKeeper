// Package app contains the application logic for the GophKeeper client.
// It includes the main application structure, its initialization, and runtime behavior.
package app

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/cli"
	"github.com/kripsy/GophKeeper/internal/client/usecase"
	"github.com/rs/zerolog"
)

// Application is the main structure of the client application.
// It holds the usecase logic and logger, orchestrating the overall behavior of the client.
type Application struct {
	usecase *usecase.ClientUsecase
	log     zerolog.Logger
}

// NewApplication creates and returns a new Application instance.
// It initializes the application with the provided configuration, build information, and logger.
// Returns an error if the initialization fails.
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

// PrepareApp prepares the application for running by setting up the user and file manager.
// Returns an error if any of the setup processes fail.
func (a *Application) PrepareApp() error {
	if err := a.usecase.SetUser(); err != nil {
		return fmt.Errorf("SetUser: %w", err)
	}

	if err := a.usecase.SetFileManager(); err != nil {
		return fmt.Errorf("SetFileManager: %w", err)
	}

	return nil
}

// Run starts the main loop of the application, allowing the user to interact with the CLI interface.
func (a *Application) Run() {
	a.usecase.InMenu()
}

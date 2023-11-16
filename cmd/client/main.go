// Package main is the entry point for the GophKeeper client application.
// It initializes and starts the client-side application,
// handling configuration, logging, and application setup.
package main

import (
	"fmt"
	"os"

	"github.com/kripsy/GophKeeper/internal/client/app"
	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/kripsy/GophKeeper/internal/client/log"
)

// main is the entry point of the application. It initializes the
// necessary configurations, logger, and the application context.
// It handles the startup sequence including error handling and
// initiates the application's run process.
func main() {
	bi := getBuildInfo()
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println("Start with config:", cfg)
	l := log.InitLogger(cfg.StoragePath)

	a, err := app.NewApplication(cfg, bi, l)
	if err != nil {
		l.Fatal().Err(err).Msg("failed create application")
	}
	if err := a.PrepareApp(); err != nil {
		fmt.Print("Something went wrong")

		return
	}

	a.Run()
}

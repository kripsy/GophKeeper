package main

import (
	"github.com/kripsy/GophKeeper/internal/client/app"
	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	bi := getBuildInfo()
	cfg := config.GetConfig()

	log := zerolog.New(os.Stdout).With().Timestamp().Logger() //todo filelog

	a, err := app.NewApplication(cfg, bi, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed create application")
	}
	a.PrepareApp()
	a.Run()
}

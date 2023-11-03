package main

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/client/app"
	"github.com/kripsy/GophKeeper/internal/client/config"
	"github.com/kripsy/GophKeeper/internal/client/log"
)

func main() {
	bi := getBuildInfo()
	cfg := config.GetConfig()
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

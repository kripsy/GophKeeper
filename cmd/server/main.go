package main

import (
	"context"
	"errors"
	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kripsy/GophKeeper/internal/logger"
	"github.com/kripsy/GophKeeper/internal/server/config"
	"github.com/kripsy/GophKeeper/internal/server/controller"
	"github.com/kripsy/GophKeeper/internal/server/infrastructure"
	"github.com/kripsy/GophKeeper/internal/server/usecase"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Printf("Error init cfg: %v", err)
		os.Exit(1)
	}

	l, err := logger.InitLog(cfg.LoggerLevel)
	if err != nil {
		fmt.Printf("Error init logger: %v", err)
		os.Exit(1)
	}

	l.Debug("Current config app: ",
		zap.String("url server", cfg.URLServer),
		zap.String("logger level", cfg.LoggerLevel),
		zap.String("database dsn", cfg.DatabaseDsn))

	repo, err := infrastructure.InitNewRepository(cfg.DatabaseDsn, l)
	if err != nil {
		l.Error("error create db instance", zap.String("msg", err.Error()))
		os.Exit(1)
	}

	userRepo, err := infrastructure.NewUserRepository(repo)
	if err != nil {
		l.Error("error init user repository", zap.String("msg", err.Error()))
		os.Exit(1)
	}
	l.Debug("NewUserRepository initialized")

	secretRepo, err := infrastructure.NewSecretRepository(repo)
	if err != nil {
		l.Error("error init secret repository", zap.String("msg", err.Error()))
		os.Exit(1)
	}
	l.Debug("NewSecretRepository initialized")

	userUseCase, err := usecase.InitUseCases(ctx, userRepo, cfg.Secret, cfg.TokenExp, l)
	if err != nil {
		l.Error("error create user usecase instance", zap.String("msg", err.Error()))
		os.Exit(1)
	}
	l.Debug("userUseCase initialized")

	secretUseCase, err := usecase.InitSecretUseCases(ctx, secretRepo, cfg.CipherSecret, l)
	if err != nil {
		l.Error("error create user usecase instance", zap.String("msg", err.Error()))
		os.Exit(1)
	}
	l.Debug("secretUseCase initialized")

	s := controller.InitNewServer(userUseCase, secretUseCase, cfg.Secret, l)

	// start shutdown application
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = s.Start(":8080")
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				l.Debug("echo stopped")
			} else {
				l.Error("error start echo", zap.String("msg", err.Error()))
			}
		}
	}()
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("Closing HTTP Server")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()

	// Закрытие репозитория и освобождение ресурсов
	l.Debug("close repository")
	if err := repo.Close(); err != nil {
		l.Error("Failed to close repository", zap.Error(err))
		os.Exit(1)
	}
	l.Debug("I'm leaving, bye!")
}

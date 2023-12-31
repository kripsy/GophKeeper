// Package main is the entry point for the GophKeeper server application.
// It sets up and starts the server, managing configurations, logging,
// database connections, and server lifecycle.
package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kripsy/GophKeeper/internal/logger"
	"github.com/kripsy/GophKeeper/internal/server/config"
	"github.com/kripsy/GophKeeper/internal/server/controller"
	"github.com/kripsy/GophKeeper/internal/server/infrastructure"
	"github.com/kripsy/GophKeeper/internal/server/usecase"
	"github.com/kripsy/GophKeeper/internal/utils"
	"go.uber.org/zap"
)

// main is the entry point of the GophKeeper server application.
// It initializes the server configuration, logging, database repositories,
// use cases, and gRPC server. It handles the application's lifecycle,
// including graceful shutdown and resource cleanup.
//
//nolint:cyclop
func main() {
	grcpPort := ":50051"
	serverCertPath := "./cert/server.crt"
	privateKeyPath := "./cert/server.key"
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Printf("Error init cfg: %v", err)

		return
	}

	l, err := logger.InitLog(cfg.LoggerLevel)
	if err != nil {
		fmt.Printf("Error init logger: %v", err)

		return
	}

	if cfg.IsSecure {
		l.Debug("creating cert")
		err = utils.CreateCertificate(serverCertPath, privateKeyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)

			return
		}
		l.Debug("cert has been created")
	}

	l.Debug("Current config app: ",
		zap.String("url server", cfg.URLServer),
		zap.String("logger level", cfg.LoggerLevel),
		zap.String("database dsn", cfg.DatabaseDsn))

	repo, err := infrastructure.InitNewRepository(cfg.DatabaseDsn, l)
	if err != nil {
		l.Error("error create db instance", zap.String("msg", err.Error()))

		return
	}

	userRepo, err := infrastructure.NewUserRepository(repo)
	if err != nil {
		l.Error("error init user repository", zap.String("msg", err.Error()))

		return
	}
	l.Debug("NewUserRepository initialized")

	l.Debug("Start init minio repository")
	ctxCreateBucket, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	minioRepo, err := infrastructure.NewMinioRepository(
		ctxCreateBucket,
		cfg.EndpointMinio,
		cfg.AccessKeyIDMinio,
		cfg.SecretAccessKeyMinio,
		cfg.BucketNameMinio,
		cfg.IsUseSSLMinio,
		l)

	if err != nil {
		l.Error("Error init minio repository", zap.Error(err))

		return
	}
	l.Debug("Success init minio repository")

	userUseCase, err := usecase.InitUseCases(ctx, userRepo, cfg.Secret, cfg.TokenExp, l)
	if err != nil {
		l.Error("error create user usecase instance", zap.String("msg", err.Error()))

		return
	}
	l.Debug("userUseCase initialized")

	secretUseCase, err := usecase.InitSecretUseCases(ctx, userRepo, minioRepo, l)
	if err != nil {
		l.Error("error create user usecase instance", zap.String("msg", err.Error()))

		return
	}
	l.Debug("secretUseCase initialized")

	s, err := controller.InitGrpcServer(userUseCase,
		secretUseCase,
		cfg.Secret,
		cfg.IsSecure,
		serverCertPath,
		privateKeyPath, l)

	if err != nil {
		l.Error("Error", zap.Error(err))

		return
	}

	// start shutdown application

	wg.Add(1)
	go func() {
		defer wg.Done()
		//nolint:gosec
		lis, err := net.Listen("tcp", grcpPort)
		if err != nil {
			l.Error("Error", zap.Error(err))

			return
		}
		l.Debug("Starting gRPC server", zap.String("port", grcpPort))
		err = s.Serve(lis)
		if err != nil {
			l.Error("Error", zap.Error(err))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("Closing GRPC Server")
		s.GracefulStop()
	}()
	wg.Wait()

	// Закрытие репозитория и освобождение ресурсов
	l.Debug("close repository")
	if err := repo.Close(); err != nil {
		l.Error("Failed to close repository", zap.Error(err))

		return
	}
	l.Debug("I'm leaving, bye!")
}

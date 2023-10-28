package controller

import (
	"context"
	"fmt"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	logger        *zap.Logger
	userUseCase   UserUseCase
	secretUseCase SecretUseCase
	secret        string
	syncStatus    *entity.SyncStatus
}

type UserUseCase interface {
	RegisterUser(ctx context.Context, user entity.User) (string, error)
	LoginUser(ctx context.Context, user entity.User) (string, error)
}

func InitGrpcServer(userUseCase UserUseCase, secretUseCase SecretUseCase, secret string, isSecure bool, logger *zap.Logger) (*grpc.Server, error) {
	syncStatus := entity.NewSyncStatus()
	s := &GrpcServer{
		userUseCase:   userUseCase,
		secretUseCase: secretUseCase,
		secret:        secret,
		logger:        logger,
		syncStatus:    syncStatus,
	}

	m := InitMyMiddleware(logger)

	interceptors := []grpc.UnaryServerInterceptor{m.AuthInterceptor}

	var srv *grpc.Server
	if isSecure {
		logger.Debug("secure grpc")
		creds, err := credentials.NewServerTLSFromFile(utils.ServerCertPath, utils.PrivateKeyPath)
		if err != nil {
			logger.Error("Failed to generate credentials.", zap.Error(err))

			return nil, fmt.Errorf("%w", err)
		}
		srv = grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	} else {
		logger.Debug("no secure grpc")
		srv = grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	}

	pb.RegisterGophKeeperServiceServer(srv, s)

	return srv, nil
}

package controller

import (
	"fmt"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const maxStreams = 20

type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	logger        *zap.Logger
	userUseCase   UserUseCase
	secretUseCase SecretUseCase
	secret        string
	syncStatus    SyncStatus
}

func InitGrpcServiceServer(userUseCase UserUseCase,
	secretUseCase SecretUseCase,
	secret string,
	logger *zap.Logger, syncStatus SyncStatus) *GrpcServer {

	return &GrpcServer{
		userUseCase:   userUseCase,
		secretUseCase: secretUseCase,
		secret:        secret,
		logger:        logger,
		syncStatus:    syncStatus,
	}
}

func InitGrpcServer(userUseCase UserUseCase,
	secretUseCase SecretUseCase,
	secret string,
	isSecure bool,
	serverCertPath string,
	privateKeyPath string,
	logger *zap.Logger) (*grpc.Server, error) {
	syncStatus := entity.NewSyncStatus()
	s := InitGrpcServiceServer(
		userUseCase,
		secretUseCase,
		secret,
		logger,
		syncStatus,
	)

	m := InitMyMiddleware(logger, secret)

	interceptors := []grpc.UnaryServerInterceptor{m.AuthInterceptor}

	var srv *grpc.Server
	if isSecure {
		logger.Debug("secure grpc")
		creds, err := credentials.NewServerTLSFromFile(serverCertPath, privateKeyPath)
		if err != nil {
			logger.Error("Failed to generate credentials.", zap.Error(err))

			return nil, fmt.Errorf("%w", err)
		}
		srv = grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)),
			grpc.MaxConcurrentStreams(maxStreams),
			grpc.StreamInterceptor(m.StreamAuthInterceptor),
		)
	} else {
		logger.Debug("no secure grpc")
		srv = grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)),
			grpc.MaxConcurrentStreams(maxStreams),
			grpc.StreamInterceptor(m.StreamAuthInterceptor),
		)
	}

	pb.RegisterGophKeeperServiceServer(srv, s)

	return srv, nil
}

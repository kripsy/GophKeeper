package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/protovalidate-go"

	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	logger        *zap.Logger
	userUseCase   UserUseCase
	secretUseCase SecretUseCase
	secret        string
}

type UserUseCase interface {
	RegisterUser(ctx context.Context, user entity.User) (string, error)
	LoginUser(ctx context.Context, user entity.User) (string, error)
}

func InitGrpcServer(userUseCase UserUseCase, secretUseCase SecretUseCase, secret string, isSecure bool, logger *zap.Logger) (*grpc.Server, error) {
	s := &GrpcServer{
		userUseCase:   userUseCase,
		secretUseCase: secretUseCase,
		secret:        secret,
		logger:        logger,
	}
	// interceptors := []grpc.UnaryServerInterceptor{}

	var srv *grpc.Server
	if isSecure {
		logger.Debug("secure grpc")
		creds, err := credentials.NewServerTLSFromFile(utils.ServerCertPath, utils.PrivateKeyPath)
		if err != nil {
			logger.Error("Failed to generate credentials.", zap.Error(err))

			return nil, fmt.Errorf("%w", err)
		}
		srv = grpc.NewServer(grpc.Creds(creds)) //grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	} else {
		logger.Debug("no secure grpc")
		srv = grpc.NewServer()
	}

	pb.RegisterGophKeeperServiceServer(srv, s)

	return srv, nil
}

func (s *GrpcServer) Register(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.logger.Debug("start Register")
	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if err = v.Validate(req); err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
	}

	user, err := entity.InitNewUser(req.Username, req.Password)
	if err != nil {
		s.logger.Error("Failed to init user", zap.Error(err))

		return nil, status.Error(codes.InvalidArgument, "Invalid request payload")
	}

	token, err := s.userUseCase.RegisterUser(ctx, user)
	if err != nil {
		var userUniqueError *models.UserExistsError
		if errors.As(err, &userUniqueError) {
			s.logger.Debug("Failed to RegisterUser user, duplicate", zap.Error(err))

			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		s.logger.Error("Failed to RegisterUser user", zap.Error(err))

		return nil, status.Error(codes.InvalidArgument, "Invalid request payload")
	}

	return &pb.AuthResponse{
		Token: token,
	}, nil
}

func (s *GrpcServer) Login(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.logger.Debug("start Login")
	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if err = v.Validate(req); err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, err.Error()))
	}

	user, err := entity.InitNewUser(req.Username, req.Password)
	if err != nil {
		s.logger.Error("Failed to init user", zap.Error(err))

		return nil, status.Error(codes.InvalidArgument, "Invalid request payload")
	}

	token, err := s.userUseCase.LoginUser(ctx, user)
	if err != nil {
		s.logger.Error("Failed to LoginUser user", zap.Error(err))

		return nil, status.Error(codes.Unauthenticated, "Failed to login")
	}

	return &pb.AuthResponse{
		Token: token,
	}, nil
}

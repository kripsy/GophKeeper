package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
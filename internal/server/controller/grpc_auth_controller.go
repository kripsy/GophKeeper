// Package controller handles gRPC server functionalities for the GophKeeper application.
// It defines the GrpcServer structure with methods for user registration, login, and server health checks.
package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register handles the user registration process via gRPC.
// It validates the request, initializes a new user, registers the user, and returns an authentication token.
func (s *GrpcServer) Register(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.logger.Debug("start Register")
	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if err = v.Validate(req); err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, "Validation error"))
	}

	user, err := entity.InitNewUser(req.GetUsername(), req.GetPassword())
	if err != nil {
		s.logger.Error("Failed to init user", zap.Error(err))

		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, "Invalid request payload"))
	}

	token, userID, err := s.userUseCase.RegisterUser(ctx, user)
	if err != nil {
		var userUniqueError *models.UserExistsError
		if errors.As(err, &userUniqueError) {
			s.logger.Debug("Failed to RegisterUser user, duplicate", zap.Error(err))

			return nil, fmt.Errorf("%w", status.Error(codes.AlreadyExists, "User already exists"))
		}
		s.logger.Error("Failed to RegisterUser user", zap.Error(err))

		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, "Invalid request payload"))
	}

	_, err = s.secretUseCase.CreateBucketSecret(ctx, user.Username, userID)
	if err != nil {
		s.logger.Debug("Failed to CreateBucketSecret", zap.Error(err))

		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	return &pb.AuthResponse{
		Token: token,
	}, nil
}

// Login handles the user login process via gRPC.
// It validates the request, authenticates the user, and returns an authentication token.
func (s *GrpcServer) Login(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.logger.Debug("start Login")
	v, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.Internal, err.Error()))
	}

	if err = v.Validate(req); err != nil {
		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, "Validation error"))
	}

	user, err := entity.InitNewUser(req.GetUsername(), req.GetPassword())
	if err != nil {
		s.logger.Error("Failed to init user", zap.Error(err))

		return nil, fmt.Errorf("%w", status.Error(codes.InvalidArgument, "Invalid request payload"))
	}

	token, userID, err := s.userUseCase.LoginUser(ctx, user)
	if err != nil {
		s.logger.Error("Failed to LoginUser user", zap.Error(err))

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "Failed to login"))
	}

	_, err = s.secretUseCase.CreateBucketSecret(ctx, user.Username, userID)
	if err != nil {
		s.logger.Debug("Failed to CreateBucketSecret", zap.Error(err))

		return nil, fmt.Errorf("%w", status.Error(codes.Internal, "failed to create bucket secret"))
	}

	return &pb.AuthResponse{
		Token: token,
	}, nil
}

// Ping is a health check method that responds to gRPC ping requests.
// It simply returns an empty response indicating that the server is operational.
func (s *GrpcServer) Ping(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.logger.Debug("start Ping")

	return &empty.Empty{}, nil
}

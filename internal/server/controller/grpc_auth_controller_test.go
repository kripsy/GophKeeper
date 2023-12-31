package controller_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/controller"
	"github.com/kripsy/GophKeeper/internal/server/controller/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcServerRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(mockCtrl)
	mockSecretUseCase := mocks.NewMockSecretUseCase(mockCtrl)

	logger := zap.NewNop()
	grpcService := controller.InitGrpcServiceServer(mockUserUseCase,
		mockSecretUseCase,
		"secret",
		logger,
		nil)

	testCases := []struct {
		name          string
		req           *pb.AuthRequest
		mockSetup     func()
		expectedToken string
		expectedErr   error
	}{
		{
			name: "Success",
			req: &pb.AuthRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return("token123", 123, nil)

				mockSecretUseCase.EXPECT().
					CreateBucketSecret(gomock.Any(), "testuser", 123).
					Return(true, nil)
			},
			expectedToken: "token123",
			expectedErr:   nil,
		},
		{
			name: "UserExistsError",
			req: &pb.AuthRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return("", 0, &models.UserExistsError{})
			},
			expectedToken: "",
			expectedErr:   status.Error(codes.AlreadyExists, "User already exists"),
		},
		{
			name: "InternalError",
			req: &pb.AuthRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return("", 0, models.NewUnionError("Invalid request payload"))
			},
			expectedToken: "",
			expectedErr:   status.Error(codes.InvalidArgument, "Invalid request payload"),
		},
		{
			name: "No valid data",
			req: &pb.AuthRequest{
				Username: "",
				Password: "testpass",
			},
			mockSetup: func() {

			},
			expectedToken: "",
			expectedErr:   status.Error(codes.InvalidArgument, "Validation error"),
		},
		{
			name: "Error create bucket",
			req: &pb.AuthRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return("token123", 123, nil)
				mockSecretUseCase.EXPECT().
					CreateBucketSecret(gomock.Any(), "testuser", 123).
					Return(false, ErrEmpty)
			},
			expectedToken: "",
			expectedErr:   status.Error(codes.Internal, "failed to create bucket secret"),
		},
		{
			name: "Failed init user",
			req: &pb.AuthRequest{
				Username: "te",
				Password: "te",
			},
			mockSetup: func() {

			},
			expectedToken: "",
			expectedErr:   status.Error(codes.InvalidArgument, "Invalid request payload"),
		},
		{
			name: "validation error",
			req: &pb.AuthRequest{
				Username: "",
				Password: "1",
			},
			mockSetup: func() {

			},
			expectedToken: "",
			expectedErr:   status.Error(codes.InvalidArgument, "Validation error"),
		},
		{
			name: "validation error",
			req: &pb.AuthRequest{
				Username: "1",
				Password: "",
			},
			mockSetup: func() {

			},
			expectedToken: "",
			expectedErr:   status.Error(codes.InvalidArgument, "Validation error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			resp, err := grpcService.Register(context.Background(), tc.req)

			if tc.expectedErr != nil {
				require.Equal(t, status.Code(tc.expectedErr), status.Code(err))
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.expectedToken, resp.GetToken())
			}
		})
	}
}

func TestGrpcServerLogin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(mockCtrl)
	mockSecretUseCase := mocks.NewMockSecretUseCase(mockCtrl)

	logger := zap.NewNop()
	grpcService := controller.InitGrpcServiceServer(mockUserUseCase,
		mockSecretUseCase,
		"secret",
		logger, nil)

	testCases := []struct {
		name          string
		req           *pb.AuthRequest
		mockSetup     func()
		expectedToken string
		expectedErr   error
	}{
		{
			name: "Success",
			req: &pb.AuthRequest{
				Username: "testuser",
				Password: "testpass",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return("token123", 123, nil)

				mockSecretUseCase.EXPECT().
					CreateBucketSecret(gomock.Any(), "testuser", 123).
					Return(true, nil)
			},
			expectedToken: "token123",
			expectedErr:   nil,
		},
		{
			name: "UnauthenticatedError",
			req: &pb.AuthRequest{
				Username: "testuser",
				Password: "wrongpass",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return("", 0, models.NewUnionError("authentication failed"))
			},
			expectedToken: "",
			expectedErr:   status.Error(codes.Unauthenticated, "Failed to login"),
		},
		{
			name: "Error create bucket",
			req: &pb.AuthRequest{
				Username: "testuser",
				Password: "wrongpass",
			},
			mockSetup: func() {
				mockUserUseCase.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return("token123", 123, nil)
				mockSecretUseCase.EXPECT().
					CreateBucketSecret(gomock.Any(), "testuser", 123).
					Return(false, ErrEmpty)
			},
			expectedToken: "",
			expectedErr:   status.Error(codes.Internal, "failed to create bucket secret"),
		},
		{
			name: "Failed init user",
			req: &pb.AuthRequest{
				Username: "te",
				Password: "te",
			},
			mockSetup: func() {

			},
			expectedToken: "",
			expectedErr:   status.Error(codes.InvalidArgument, "Invalid request payload"),
		},
		{
			name: "validation error",
			req: &pb.AuthRequest{
				Username: "",
				Password: "1",
			},
			mockSetup: func() {

			},
			expectedToken: "",
			expectedErr:   status.Error(codes.InvalidArgument, "Validation error"),
		},
		{
			name: "validation error",
			req: &pb.AuthRequest{
				Username: "1",
				Password: "",
			},
			mockSetup: func() {

			},
			expectedToken: "",
			expectedErr:   status.Error(codes.InvalidArgument, "Validation error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			resp, err := grpcService.Login(context.Background(), tc.req)

			if tc.expectedErr != nil {
				assert.Equal(t, status.Code(tc.expectedErr), status.Code(err))
				fmt.Println(err.Error())
				require.EqualError(t, tc.expectedErr, err.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.expectedToken, resp.GetToken())
			}
		})
	}
}

func TestGrpcServerPing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zap.NewNop()
	grpcService := controller.InitGrpcServiceServer(nil, nil, "secret", logger, nil)

	testCases := []struct {
		name        string
		req         *empty.Empty
		expectedErr error
	}{
		{
			name:        "Success",
			req:         &empty.Empty{},
			expectedErr: nil,
		},
		// Так как функция Ping не генерирует ошибок, нет необходимости в дополнительных тестовых случаях.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := grpcService.Ping(context.Background(), tc.req)

			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}

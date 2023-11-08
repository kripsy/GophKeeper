package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/GophKeeper/internal/server/controller"
	"github.com/kripsy/GophKeeper/internal/server/controller/mocks"
	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestInitMyMiddleware(t *testing.T) {
	mockLogger := zap.NewNop() // Используем No-op Logger от zap для тестирования

	tests := []struct {
		name   string
		secret string
	}{
		{
			name:   "Success",
			secret: "secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Переинициализация одиночки перед каждым тестом

			middleware := controller.InitMyMiddleware(mockLogger, tt.secret)
			require.NotNil(t, middleware)

			// Проверяем, что повторный вызов возвращает тот же экземпляр
			sameMiddleware := controller.InitMyMiddleware(mockLogger, "newSecret")
			require.Equal(t, middleware, sameMiddleware)
		})
	}
}

func TestAuthInterceptor(t *testing.T) {
	//nolint:goerr113
	var ErrContext = errors.New("context error")
	const (
		loginMethod    = "/pkg.api.gophkeeper.v1.GophKeeperService/Login"
		registerMethod = "/pkg.api.gophkeeper.v1.GophKeeperService/Register"
		pingMethod     = "/pkg.api.gophkeeper.v1.GophKeeperService/Ping"
	)
	l := zap.NewNop()
	secret := "secret"

	middleware := controller.InitMyMiddleware(l, secret)

	userID := 1
	userName := "testuser"
	validToken, err := utils.BuildJWTString(userID, userName, secret, time.Hour)
	require.NoError(t, err)

	invalidToken, err := utils.BuildJWTString(userID, userName, secret+"fake", time.Hour)
	require.NoError(t, err)

	tests := []struct {
		name        string
		fullMethod  string
		token       string
		secret      string
		wantErr     bool
		isProtected bool
	}{
		{
			name:        "ValidToken",
			fullMethod:  "/pkg.api.gophkeeper.v1.GophKeeperService/ValidMethod",
			token:       validToken,
			secret:      secret,
			wantErr:     false,
			isProtected: true,
		},
		{
			name:        "InvalidToken",
			fullMethod:  "/pkg.api.gophkeeper.v1.GophKeeperService/ValidMethod",
			token:       invalidToken,
			secret:      secret,
			wantErr:     true,
			isProtected: true,
		},
		{
			name:        "NoMetadata",
			fullMethod:  "/pkg.api.gophkeeper.v1.GophKeeperService/ValidMethod",
			token:       "",
			secret:      secret,
			wantErr:     true,
			isProtected: true,
		},
		{
			name:        "NoAuthHeader",
			fullMethod:  "/pkg.api.gophkeeper.v1.GophKeeperService/ValidMethod",
			token:       "",
			secret:      secret,
			wantErr:     true,
			isProtected: true,
		},
		{
			name:        "EmptyToken",
			fullMethod:  "/pkg.api.gophkeeper.v1.GophKeeperService/ValidMethod",
			token:       "Bearer ",
			secret:      secret,
			wantErr:     true,
			isProtected: true,
		},
		{
			name:        "UnprotectedMethod",
			fullMethod:  loginMethod,
			token:       validToken,
			secret:      secret,
			wantErr:     false,
			isProtected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				if tt.isProtected {
					if ctx.Value(utils.USERNAMECONTEXTKEY) != userName {
						return "", ErrContext
					}
					if ctx.Value(utils.USERIDCONTEXTKEY) != userID {
						return "", ErrContext
					}
				}

				return "testResponse", nil
			}

			md := metadata.New(map[string]string{
				utils.AUTHORIZATIONHEADER: tt.token,
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			switch tt.name {
			case "NoMetadata":
				ctx = context.Background()
			case "NoAuthHeader":
				md := metadata.New(map[string]string{})
				ctx = metadata.NewIncomingContext(context.Background(), md)
			case "EmptyToken":
				md := metadata.New(map[string]string{
					utils.AUTHORIZATIONHEADER: "Bearer ",
				})
				ctx = metadata.NewIncomingContext(context.Background(), md)
			}

			info := &grpc.UnaryServerInfo{
				FullMethod: tt.fullMethod,
			}

			_, err := middleware.AuthInterceptor(ctx, nil, info, testHandler)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStreamAuthInterceptor(t *testing.T) {
	secret := "secret"
	logger := zap.NewNop()
	middleware := controller.InitMyMiddleware(logger, secret)
	userID := 1
	userName := "testuser"
	validToken, _ := utils.BuildJWTString(userID, userName, secret, time.Hour)
	invalidToken, _ := utils.BuildJWTString(userID, userName, "fake", time.Hour)

	tests := []struct {
		name       string
		fullMethod string
		token      string
		wantErr    bool
		wantCode   codes.Code
		setup      func(*mocks.MockServerStream)
	}{
		{
			name:       "ValidToken",
			fullMethod: "/package.Service/StreamMethod",
			token:      validToken,
			wantErr:    false,
			setup: func(mockStream *mocks.MockServerStream) {
				md := metadata.New(map[string]string{
					utils.AUTHORIZATIONHEADER: utils.TOKENPREFIX + validToken,
				})
				ctx := metadata.NewIncomingContext(context.Background(), md)
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
		},
		{
			name:       "InvalidToken",
			fullMethod: "/package.Service/StreamMethod",
			token:      invalidToken,
			wantErr:    true,
			wantCode:   codes.Unauthenticated,
			setup: func(mockStream *mocks.MockServerStream) {
				md := metadata.New(map[string]string{
					utils.AUTHORIZATIONHEADER: utils.TOKENPREFIX + invalidToken,
				})
				ctx := metadata.NewIncomingContext(context.Background(), md)
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
		},
		{
			name:       "MetadataNotFound",
			fullMethod: "/package.Service/StreamMethod",
			token:      "",
			wantErr:    true,
			wantCode:   codes.Internal,
			setup: func(mockStream *mocks.MockServerStream) {
				ctx := context.Background()
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
		},
		{
			name:       "AuthorizationHeaderMissing",
			fullMethod: "/package.Service/StreamMethod",
			token:      "",
			wantErr:    true,
			setup: func(mockStream *mocks.MockServerStream) {
				md := metadata.New(map[string]string{})
				ctx := metadata.NewIncomingContext(context.Background(), md)
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStream := mocks.NewMockServerStream(ctrl)
			tt.setup(mockStream)

			info := &grpc.StreamServerInfo{
				FullMethod: tt.fullMethod,
			}

			err := middleware.StreamAuthInterceptor(nil, mockStream,
				info, func(srv interface{}, stream grpc.ServerStream) error {
					return nil
				})

			if tt.wantErr {
				require.Error(t, err)
				if err != nil {
					require.Equal(t, tt.wantCode, status.Code(err))
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

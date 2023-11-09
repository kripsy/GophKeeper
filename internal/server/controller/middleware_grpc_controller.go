package controller

import (
	"context"
	"fmt"
	"strings"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kripsy/GophKeeper/internal/utils/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MyMiddleware struct {
	myLogger *zap.Logger
	secret   string
}

var (
	//nolint:gochecknoglobals
	once sync.Once
	//nolint:gochecknoglobals
	instance *MyMiddleware
)

const (
	loginMethod    = "/pkg.api.gophkeeper.v1.GophKeeperService/Login"
	registerMethod = "/pkg.api.gophkeeper.v1.GophKeeperService/Register"
	pingMethod     = "/pkg.api.gophkeeper.v1.GophKeeperService/Ping"
)

func InitMyMiddleware(myLogger *zap.Logger, secret string) *MyMiddleware {
	once.Do(func() {
		instance = &MyMiddleware{
			myLogger: myLogger,
			secret:   secret,
		}
	})

	return instance
}

func (m MyMiddleware) AuthInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	case loginMethod,
		registerMethod,
		pingMethod:
		m.myLogger.Debug("No protected method", zap.String("method", info.FullMethod))

		return handler(ctx, req)
	}

	m.myLogger.Debug("Protected method")
	m.myLogger.Debug(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		m.myLogger.Debug("couldn't extract metadata from req")

		return nil, fmt.Errorf("%w", status.Error(codes.Internal, "couldn't extract metadata from req"))
	}

	authHeaders, ok := md[auth.AUTHORIZATIONHEADER]
	if !ok || len(authHeaders) != 1 {
		m.myLogger.Debug("authorization not exists")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "authorization not exists"))
	}

	token := strings.TrimPrefix(authHeaders[0], auth.TOKENPREFIX)
	if token == "" {
		m.myLogger.Debug("token empty or not valid")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	if isValid, err := auth.IsValidToken(token, m.secret); err != nil || !isValid {
		m.myLogger.Debug("token is not valid")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}
	username, err := auth.GetUsernameFromToken(token, m.secret)
	if err != nil {
		m.myLogger.Debug("cannot get username")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	userID, err := auth.GetUseIDFromToken(token, m.secret)
	if err != nil {
		m.myLogger.Debug("cannot get userID")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	//nolint:staticcheck
	newCtx := context.WithValue(ctx, auth.USERNAMECONTEXTKEY, username)
	//nolint:staticcheck
	newCtx = context.WithValue(newCtx, auth.USERIDCONTEXTKEY, userID)

	return handler(newCtx, req)
}

func (m MyMiddleware) StreamAuthInterceptor(srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	// Получаем контекст из потокового сервера
	ctx := ss.Context()

	m.myLogger.Debug("Protected method")
	m.myLogger.Debug(info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		m.myLogger.Debug("couldn't extract metadata from req")

		return fmt.Errorf("%w", status.Error(codes.Internal, "couldn't extract metadata from req"))
	}

	authHeaders, ok := md[auth.AUTHORIZATIONHEADER]
	if !ok || len(authHeaders) != 1 {
		m.myLogger.Debug("authorization not exists")

		return fmt.Errorf("%w", status.Error(codes.Unauthenticated, "authorization not exists"))
	}

	token := strings.TrimPrefix(authHeaders[0], auth.TOKENPREFIX)
	if token == "" {
		m.myLogger.Debug("token empty or not valid")

		return fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	if isValid, err := auth.IsValidToken(token, m.secret); err != nil || !isValid {
		m.myLogger.Debug("token is not valid")

		return fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	username, err := auth.GetUsernameFromToken(token, m.secret)
	if err != nil {
		m.myLogger.Debug("cannot get username")

		return fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	userID, err := auth.GetUseIDFromToken(token, m.secret)
	if err != nil {
		m.myLogger.Debug("cannot get userID")

		return fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	// Добавляем имя пользователя в контекст
	//nolint:staticcheck
	newCtx := context.WithValue(ctx, auth.USERNAMECONTEXTKEY, username)
	// Заменяем старый контекст новым в потоковом сервере
	wrappedStream := grpc_middleware.WrapServerStream(ss)
	wrappedStream.WrappedContext = newCtx

	//nolint:staticcheck
	newCtx = context.WithValue(newCtx, auth.USERIDCONTEXTKEY, userID)
	wrappedStream = grpc_middleware.WrapServerStream(ss)
	wrappedStream.WrappedContext = newCtx

	return handler(srv, wrappedStream)
}

package controller

import (
	"context"
	"strings"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kripsy/GophKeeper/internal/utils"
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

func InitMyMiddleware(myLogger *zap.Logger, secret string) *MyMiddleware {
	once.Do(func() {
		instance = &MyMiddleware{
			myLogger: myLogger,
			secret:   secret,
		}
	})

	return instance
}

func (m MyMiddleware) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	case "/pkg.api.gophkeeper.v1.GophKeeperService/Login", "/pkg.api.gophkeeper.v1.GophKeeperService/Register":
		m.myLogger.Debug("No protected method")

		return handler(ctx, req)
	}

	m.myLogger.Debug("Protected method")
	m.myLogger.Debug(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		m.myLogger.Debug("couldn't extract metadata from req")
		return nil, status.Error(codes.Internal, "couldn't extract metadata from req")
	}

	authHeaders, ok := md[utils.AUTHORIZATIONHEADER]
	if !ok || len(authHeaders) != 1 {
		m.myLogger.Debug("authorization not exists")
		return nil, status.Error(codes.Unauthenticated, "authorization not exists")
	}

	token := strings.TrimPrefix(authHeaders[0], utils.TOKENPREFIX)
	if token == "" {
		m.myLogger.Debug("token empty or not valid")
		return nil, status.Error(codes.Unauthenticated, "token empty or not valid")
	}

	if isValid, err := utils.IsValidToken(token, m.secret); err != nil || !isValid {
		m.myLogger.Debug("token is not valid")
		return nil, status.Error(codes.Unauthenticated, "token empty or not valid")
	}
	username, err := utils.GetUsernameFromToken(token, m.secret)
	if err != nil {
		m.myLogger.Debug("cannot get username")
		return nil, status.Error(codes.Unauthenticated, "token empty or not valid")
	}

	newCtx := context.WithValue(ctx, utils.USERNAMECONTEXTKEY, username)

	return handler(newCtx, req)
}

func (m MyMiddleware) StreamAuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// Получаем контекст из потокового сервера
	ctx := ss.Context()

	switch info.FullMethod {
	case "/pkg.api.gophkeeper.v1.GophKeeperService/LoginStream", "/pkg.api.gophkeeper.v1.GophKeeperService/RegisterStream":
		m.myLogger.Debug("No protected method")
		return handler(srv, ss)
	}

	m.myLogger.Debug("Protected method")
	m.myLogger.Debug(info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		m.myLogger.Debug("couldn't extract metadata from req")
		return status.Error(codes.Internal, "couldn't extract metadata from req")
	}

	authHeaders, ok := md[utils.AUTHORIZATIONHEADER]
	if !ok || len(authHeaders) != 1 {
		m.myLogger.Debug("authorization not exists")
		return status.Error(codes.Unauthenticated, "authorization not exists")
	}

	token := strings.TrimPrefix(authHeaders[0], utils.TOKENPREFIX)
	if token == "" {
		m.myLogger.Debug("token empty or not valid")
		return status.Error(codes.Unauthenticated, "token empty or not valid")
	}

	if isValid, err := utils.IsValidToken(token, m.secret); err != nil || !isValid {
		m.myLogger.Debug("token is not valid")
		return status.Error(codes.Unauthenticated, "token empty or not valid")
	}

	username, err := utils.GetUsernameFromToken(token, m.secret)
	if err != nil {
		m.myLogger.Debug("cannot get username")
		return status.Error(codes.Unauthenticated, "token empty or not valid")
	}

	userID, err := utils.GetUseIDFromToken(token, m.secret)
	if err != nil {
		m.myLogger.Debug("cannot get userID")
		return status.Error(codes.Unauthenticated, "token empty or not valid")
	}

	// Добавляем имя пользователя в контекст
	newCtx := context.WithValue(ctx, utils.USERNAMECONTEXTKEY, username)
	// Заменяем старый контекст новым в потоковом сервере
	wrappedStream := grpc_middleware.WrapServerStream(ss)
	wrappedStream.WrappedContext = newCtx

	newCtx = context.WithValue(newCtx, utils.USERIDCONTEXTKEY, userID)
	wrappedStream = grpc_middleware.WrapServerStream(ss)
	wrappedStream.WrappedContext = newCtx

	return handler(srv, wrappedStream)
}

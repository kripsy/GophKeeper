package controller

import (
	"context"
	"strings"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	tokenPrefix         = ""
	tokenContextKey     = "token"
)

type MyMiddleware struct {
	MyLogger *zap.Logger
}

var (
	//nolint:gochecknoglobals
	once sync.Once
	//nolint:gochecknoglobals
	instance *MyMiddleware
)

func InitMyMiddleware(myLogger *zap.Logger) *MyMiddleware {
	once.Do(func() {
		instance = &MyMiddleware{
			MyLogger: myLogger,
		}
	})

	return instance
}

func (m MyMiddleware) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	case "/pkg.api.gophkeeper.v1.GophKeeperService/Login", "/pkg.api.gophkeeper.v1.GophKeeperService/Register":
		m.MyLogger.Debug("No protected method")

		return handler(ctx, req)
	}

	m.MyLogger.Debug("Protected method")
	m.MyLogger.Debug(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		m.MyLogger.Debug("couldn't extract metadata from req")
		return nil, status.Error(codes.Internal, "couldn't extract metadata from req")
	}

	authHeaders, ok := md[authorizationHeader]
	if !ok || len(authHeaders) != 1 {
		m.MyLogger.Debug("authorization not exists")
		return nil, status.Error(codes.Unauthenticated, "authorization not exists")
	}

	token := strings.TrimPrefix(authHeaders[0], tokenPrefix)
	if token == "" {
		m.MyLogger.Debug("token empty or not valid")
		return nil, status.Error(codes.Unauthenticated, "token empty or not valid")
	}

	newCtx := context.WithValue(ctx, tokenContextKey, token)

	return handler(newCtx, req)
}

func ExtractTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenContextKey).(string)
	return token, ok
}

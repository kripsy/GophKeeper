package utils

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}

func AddJwtIntoContext(ctx context.Context, jwt string) context.Context {
	// Создание метаданных и добавление заголовка авторизации
	md := metadata.New(map[string]string{
		"authorization": jwt,
	})

	// Прикрепление метаданных к контексту
	newCtx := metadata.NewOutgoingContext(ctx, md)
	return newCtx
}

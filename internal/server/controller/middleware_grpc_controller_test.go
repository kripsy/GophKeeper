package controller_test

import (
	"testing"

	"github.com/kripsy/GophKeeper/internal/server/controller"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
			assert.NotNil(t, middleware)

			// Проверяем, что повторный вызов возвращает тот же экземпляр
			sameMiddleware := controller.InitMyMiddleware(mockLogger, "newSecret")
			assert.Equal(t, middleware, sameMiddleware)
		})
	}
}

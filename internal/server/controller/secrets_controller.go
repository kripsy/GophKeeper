package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type SecretPayload struct {
	Type        string                 `json:"type" validate:"required,oneof=text binary card login_password"`
	Data        map[string]interface{} `json:"data" validate:"required"`
	Meta        string                 `json:"meta,omitempty"`
	ChunkNum    int                    `json:"chunk_num,omitempty"`
	TotalChunks int                    `json:"total_chunks,omitempty"`
}

type SecretUseCase interface {
	SaveTextSecret(ctx context.Context, secret entity.Secret) (int, error)
	GetSecret(ctx context.Context, secretID int) (entity.Secret, error)
	GetSecretsByUserID(ctx context.Context, userID, limit, offset int) ([]entity.Secret, error)
}

func (s *Server) createSecret(l *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user_id").(int) // Получаем токен из контекста
		var payload SecretPayload
		if err := c.Bind(&payload); err != nil {
			l.Error("Failed to bind request body", zap.Error(err))
			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Invalid request payload",
			})
		}

		if err := c.Validate(&payload); err != nil {
			l.Error("Validation error", zap.Error(err))
			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Validation error",
			})
		}

		ctx := c.Request().Context()

		// Определение, какой usecase вызывать на основе значения Type
		switch payload.Type {
		case "text":
			// Вызов usecase для текстовых секретов
			l.Debug("income text secret")
			data, err := json.Marshal(payload.Data)
			if err != nil {
				l.Error("Failed to marshall secret", zap.Error(err))

				return c.JSON(http.StatusInternalServerError, Response{
					Data:  nil,
					Error: "Failed to marshall secret",
				})
			}
			s.secretUseCase.SaveTextSecret(ctx, entity.Secret{
				Type:        payload.Type,
				Data:        data,
				Meta:        payload.Meta,
				ChunkNum:    payload.ChunkNum,
				TotalChunks: payload.TotalChunks,
				UserID:      userID,
			})
		case "binary":
			// Вызов usecase для бинарных секретов
			l.Debug("income binary")
		case "card":
			// Вызов usecase для секретов типа карта
			l.Debug("income card")
		case "login_password":
			// Вызов usecase для секретов типа логин/пароль
			l.Debug("income login_password")
		default:
			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Unknown secret type",
			})
		}

		// Возвращаем успешный ответ
		return c.JSON(http.StatusOK, Response{
			Data:  map[string]string{"message": "Secret successfully saved"},
			Error: nil,
		})
	}
}

func (s *Server) getSecretsByUserID(l *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user_id").(int)

		limitParam := c.QueryParam("limit")
		offsetParam := c.QueryParam("offset")

		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			limit = 10
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			offset = 0
		}

		ok, err := utils.CheckValidQueryParam(limit, offset)
		if err != nil {
			l.Error("Error check offset, limit")

			return c.JSON(http.StatusInternalServerError, Response{
				Data:  nil,
				Error: "Error check offset, limit",
			})
		}

		if !ok {
			l.Error("Incorrect offset, limit")

			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Incorrect offset, limit",
			})
		}

		secrets, err := s.secretUseCase.GetSecretsByUserID(c.Request().Context(), userID, limit, offset)
		if err != nil {
			l.Error("Failed to get secrets", zap.Error(err))

			return c.JSON(http.StatusInternalServerError, Response{
				Data:  nil,
				Error: "Failed to get secrets",
			})
		}
		if len(secrets) < 1 {
			return c.JSON(http.StatusNoContent, Response{
				Data:  nil,
				Error: nil,
			})
		}

		return c.JSON(http.StatusOK, Response{
			Data:  secrets,
			Error: nil,
		})
	}
}

func (s *Server) getSecret(l *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		secretID := 1 //c.Param("id") // Assuming you're getting the secret by ID from the URL

		ctx := c.Request().Context()
		secret, err := s.secretUseCase.GetSecret(ctx, secretID)

		if err != nil {
			l.Error("Failed to GetSecret", zap.Error(err))

			return c.JSON(http.StatusNotFound, Response{
				Data:  nil,
				Error: "Secret not found",
			})
		}

		return c.JSON(http.StatusOK, Response{
			Data:  secret,
			Error: nil,
		})
	}
}

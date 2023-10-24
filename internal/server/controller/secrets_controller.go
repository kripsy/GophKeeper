package controller

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/kripsy/GophKeeper/internal/models"
// 	"github.com/kripsy/GophKeeper/internal/server/entity"
// 	"github.com/kripsy/GophKeeper/internal/utils"
// 	"github.com/labstack/echo/v4"
// 	"go.uber.org/zap"
// )

// type SecretPayload struct {
// 	Type        string                 `json:"type" validate:"required,oneof=text binary card login_password"`
// 	Data        map[string]interface{} `json:"data" validate:"required"`
// 	Meta        string                 `json:"meta,omitempty"`
// 	ChunkNum    int                    `json:"chunk_num,omitempty"`
// 	TotalChunks int                    `json:"total_chunks,omitempty"`
// }

// type TextSecret struct {
// 	Content string `json:"content" validate:"required"`
// }

// func (s *Server) createSecret(l *zap.Logger) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		userID := c.Get("user_id").(int) // Получаем токен из контекста
// 		var payload SecretPayload
// 		if err := c.Bind(&payload); err != nil {
// 			l.Error("Failed to bind request body", zap.Error(err))
// 			return c.JSON(http.StatusBadRequest, Response{
// 				Data:  nil,
// 				Error: "Invalid request payload",
// 			})
// 		}

// 		if err := c.Validate(&payload); err != nil {
// 			l.Error("Validation error", zap.Error(err))
// 			return c.JSON(http.StatusBadRequest, Response{
// 				Data:  nil,
// 				Error: "Validation error",
// 			})
// 		}

// 		ctx := c.Request().Context()

// 		// Определение, какой usecase вызывать на основе значения Type
// 		switch payload.Type {
// 		case "text":
// 			// Вызов usecase для текстовых секретов
// 			l.Debug("income text secret")
// 			data, err := json.Marshal(payload.Data)
// 			if err != nil {
// 				l.Error("Failed to marshall secret", zap.Error(err))

// 				return c.JSON(http.StatusInternalServerError, Response{
// 					Data:  nil,
// 					Error: "Failed to marshall secret",
// 				})
// 			}
// 			s.secretUseCase.SaveTextSecret(ctx, entity.Secret{
// 				Type:        payload.Type,
// 				Data:        data,
// 				Meta:        payload.Meta,
// 				ChunkNum:    payload.ChunkNum,
// 				TotalChunks: payload.TotalChunks,
// 				UserID:      userID,
// 			})
// 		case "binary":
// 			// Вызов usecase для бинарных секретов
// 			l.Debug("income binary")
// 		case "card":
// 			// Вызов usecase для секретов типа карта
// 			l.Debug("income card")
// 		case "login_password":
// 			// Вызов usecase для секретов типа логин/пароль
// 			l.Debug("income login_password")
// 		default:
// 			return c.JSON(http.StatusBadRequest, Response{
// 				Data:  nil,
// 				Error: "Unknown secret type",
// 			})
// 		}

// 		// Возвращаем успешный ответ
// 		return c.JSON(http.StatusOK, Response{
// 			Data:  map[string]string{"message": "Secret successfully saved"},
// 			Error: nil,
// 		})
// 	}
// }

// func (s *Server) getSecretsByUserID(l *zap.Logger) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		userID := c.Get("user_id").(int)

// 		limitParam := c.QueryParam("limit")
// 		offsetParam := c.QueryParam("offset")

// 		limit, err := strconv.Atoi(limitParam)
// 		if err != nil {
// 			limit = 10
// 		}

// 		offset, err := strconv.Atoi(offsetParam)
// 		if err != nil {
// 			offset = 0
// 		}

// 		ok, err := utils.CheckValidQueryParam(limit, offset)
// 		if err != nil {
// 			l.Error("Error check offset, limit")

// 			return c.JSON(http.StatusInternalServerError, Response{
// 				Data:  nil,
// 				Error: "Error check offset, limit",
// 			})
// 		}

// 		if !ok {
// 			l.Error("Incorrect offset, limit")

// 			return c.JSON(http.StatusBadRequest, Response{
// 				Data:  nil,
// 				Error: "Incorrect offset, limit",
// 			})
// 		}

// 		secrets, err := s.secretUseCase.GetSecretsByUserID(c.Request().Context(), userID, limit, offset)
// 		if err != nil {
// 			l.Error("Failed to get secrets", zap.Error(err))

// 			return c.JSON(http.StatusInternalServerError, Response{
// 				Data:  nil,
// 				Error: "Failed to get secrets",
// 			})
// 		}
// 		if len(secrets) < 1 {
// 			return c.JSON(http.StatusNoContent, Response{
// 				Data:  nil,
// 				Error: nil,
// 			})
// 		}

// 		return c.JSON(http.StatusOK, Response{
// 			Data:  secrets,
// 			Error: nil,
// 		})
// 	}
// }

// func (s *Server) getSecretByID(l *zap.Logger) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		ctx := c.Request().Context()
// 		// Получение ID секрета из URL
// 		secretIDStr := c.Param("id")
// 		secretID, err := strconv.Atoi(secretIDStr)
// 		if err != nil {
// 			l.Error("Invalid secret ID", zap.Error(err))

// 			return c.JSON(http.StatusBadRequest, Response{
// 				Data:  nil,
// 				Error: "Invalid secret ID",
// 			})
// 		}

// 		// Извлечение userID из контекста (предполагая, что он там есть после аутентификации)
// 		userID, ok := c.Get("user_id").(int)
// 		if !ok {
// 			l.Error("Failed to retrieve userID from context")

// 			return c.JSON(http.StatusInternalServerError, Response{
// 				Data:  nil,
// 				Error: "Failed to retrieve userID",
// 			})
// 		}

// 		// Вызов usecase для получения секрета по ID и userID
// 		secret, err := s.secretUseCase.GetSecretByID(ctx, secretID, userID)
// 		if err != nil {
// 			var secretNotFoundError *models.SecretNotFoundError
// 			if errors.As(err, &secretNotFoundError) {
// 				l.Error("Secret not found", zap.Int("userID", userID), zap.Int("secretID", secretID))

// 				return c.JSON(http.StatusNoContent, Response{
// 					Data:  nil,
// 					Error: nil,
// 				})
// 			}

// 			l.Error("Failed to retrieve secret", zap.Error(err))

// 			return c.JSON(http.StatusInternalServerError, Response{
// 				Data:  nil,
// 				Error: "Failed to retrieve secret",
// 			})
// 		}
// 		payload, err := convertSecret2Payload(secret)
// 		if err != nil {
// 			l.Error("Failed convertSecret2Payload", zap.Error(err))

// 			return c.JSON(http.StatusInternalServerError, Response{
// 				Data:  nil,
// 				Error: "Failed to retrieve secret",
// 			})
// 		}

// 		// Возвращаем секрет в ответе
// 		return c.JSON(http.StatusOK, Response{
// 			Data:  payload,
// 			Error: nil,
// 		})
// 	}
// }

// func convertSecret2Payload(secret entity.Secret) (SecretPayload, error) {
// 	payload := SecretPayload{
// 		Type:        secret.Type,
// 		Meta:        secret.Meta,
// 		ChunkNum:    secret.ChunkNum,
// 		TotalChunks: secret.TotalChunks,
// 		Data:        make(map[string]interface{}),
// 	}
// 	switch secret.Type {
// 	case "text":
// 		var textData TextSecret
// 		fmt.Println(string(secret.Data))
// 		err := json.Unmarshal(secret.Data, &textData)
// 		if err != nil {
// 			return SecretPayload{}, fmt.Errorf("%w", err)
// 		}
// 		payload.Data["content"] = textData.Content

// 	case "binary":

// 	case "card":

// 	case "login_password":
// 	}

// 	return payload, nil
// }

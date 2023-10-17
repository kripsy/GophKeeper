package controller

import (
	"io"
	"net/http"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) createSecret(l *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Data:  nil,
				Error: "Failed to read request body",
			})
		}
		defer c.Request().Body.Close()
		_, err = entity.InitNewSecret(body) // Assuming you have a function to initialize a new secret
		if err != nil {
			l.Error("Failed to init secret", zap.Error(err))

			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Invalid request payload",
			})
		}
		ctx := c.Request().Context()
		// /secretID, err := s.secretUseCase.SaveSecret(ctx, entity.User{})
		err = s.secretUseCase.SaveSecret(ctx, entity.User{})

		if err != nil {
			l.Error("Failed to CreateSecret", zap.Error(err))

			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Failed to create secret",
			})
		}

		responseData := make(map[string]int)
		responseData["secret_id"] = 1 // secretID

		return c.JSON(http.StatusOK, Response{
			Data:  responseData,
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

// You can add more functions like updateSecret, deleteSecret, etc.

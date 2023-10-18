package controller

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserUseCase interface {
	RegisterUser(ctx context.Context, user entity.User) (string, error)
	LoginUser(ctx context.Context, user entity.User) (string, error)
}

func (s *Server) register(l *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Data:  nil,
				Error: "Failed to read request body",
			})
		}
		defer c.Request().Body.Close()
		user, err := entity.InitNewUser(body)
		if err != nil {
			l.Error("Failed to init user", zap.Error(err))

			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Invalid request payload",
			})
		}
		ctx := c.Request().Context()
		token, err := s.userUseCase.RegisterUser(ctx, user)

		if err != nil {
			l.Error("Failed to RegisterUser user", zap.Error(err))
			var userUniqueError *models.UserExistsError
			if errors.As(err, &userUniqueError) {
				return c.JSON(http.StatusConflict, Response{
					Data:  nil,
					Error: "User already exists",
				})
			}

			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Invalid request payload",
			})
		}
		responseData := make(map[string]string)
		responseData["token"] = token

		return c.JSON(http.StatusOK, Response{
			Data:  responseData,
			Error: nil,
		})
	}
}

func (s *Server) login(l *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Data:  nil,
				Error: "Failed to read request body",
			})
		}
		defer c.Request().Body.Close()
		user, err := entity.InitNewUser(body)
		if err != nil {
			l.Error("Failed to init user", zap.Error(err))

			return c.JSON(http.StatusBadRequest, Response{
				Data:  nil,
				Error: "Invalid request payload",
			})
		}
		ctx := c.Request().Context()
		token, err := s.userUseCase.LoginUser(ctx, user)

		if err != nil {
			l.Error("Failed to LoginUser user", zap.Error(err))

			return c.JSON(http.StatusUnauthorized, Response{
				Data:  nil,
				Error: "Failed to login",
			})
		}

		responseData := make(map[string]string)
		responseData["token"] = token

		return c.JSON(http.StatusOK, Response{
			Data:  responseData,
			Error: nil,
		})
	}
}

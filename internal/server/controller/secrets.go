package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) addSecret(l *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.JSON(http.StatusOK)
	}
}

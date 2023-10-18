package controller

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.uber.org/zap"
)

const CONTEXTTIMEOUT = 2 * time.Second

type Response struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type Server struct {
	server        *echo.Echo
	userUseCase   UserUseCase
	secretUseCase SecretUseCase
	secret        string
	// middleware  *deviceHandlerMiddleware
}

func InitNewServer(userUseCase UserUseCase, secretUseCase SecretUseCase, secret string, logger *zap.Logger) *Server {
	s := &Server{
		server:        echo.New(),
		userUseCase:   userUseCase,
		secretUseCase: secretUseCase,
		secret:        secret,
		// middleware:  InitDeviceHandlerMiddleware(logger),
	}
	s.setRouting(logger)
	return s
}

func (s *Server) Start(addr string) error {
	return s.server.Start(addr)
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) setRouting(logger *zap.Logger) error {
	// s.server.Use(middleware.CORS())
	s.server.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: CONTEXTTIMEOUT,
	}))
	s.server.Validator = &CustomValidator{validator: validator.New()}

	authGroup := s.server.Group("/users")
	authGroup.POST("/register", s.register(logger))
	authGroup.POST("/login", s.login(logger))

	secretGroup := s.server.Group("/secrets")
	secretGroup.Use(TokenAuthMiddleware(logger, s.secret))
	secretGroup.GET("", s.getSecretsByUserID(logger))
	secretGroup.GET("/:id", s.getSecretByID(logger))
	secretGroup.POST("/create", s.createSecret(logger))

	// // create routes group
	// devicesGroup := s.server.Group("/devices")
	// // apply middleware for group
	// devicesGroup.Use(s.middleware.DeviceHandlerMiddleware)

	// devicesGroup.GET("/:id", s.GetDeviceEventsByID(logger))
	// devicesGroup.POST("/location", s.SetDeviceLocations(logger))
	// devicesGroup.GET("", s.GetDevices(logger))

	return nil
}

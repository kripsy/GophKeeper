package controller

import (
	"context"

	"github.com/kripsy/GophKeeper/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.uber.org/zap"
)

type Server struct {
	server      *echo.Echo
	userUseCase usecase.UserUseCase
	// middleware  *deviceHandlerMiddleware
}

func InitNewServer(userUseCase usecase.UserUseCase, logger *zap.Logger) *Server {

	s := &Server{
		server:      echo.New(),
		userUseCase: userUseCase,
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
	s.server.Use(middleware.CORS())

	authGroup := s.server.Group("/users")
	authGroup.POST("/register", s.login(logger))
	authGroup.POST("/login", s.login(logger))

	// // create routes group
	// devicesGroup := s.server.Group("/devices")
	// // apply middleware for group
	// devicesGroup.Use(s.middleware.DeviceHandlerMiddleware)

	// devicesGroup.GET("/:id", s.GetDeviceEventsByID(logger))
	// devicesGroup.POST("/location", s.SetDeviceLocations(logger))
	// devicesGroup.GET("", s.GetDevices(logger))

	return nil
}

func (s *Server) login(l *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(200, "not implemented")
	}
}

// func (s *Server) SetDeviceLocations(logger *zap.Logger) echo.HandlerFunc {
// 	return func(ctx echo.Context) error {
// 		var deviceEventRequest entity.DeviceEventRequest
// 		deviceEventRequest, ok := ctx.Get("deviceEventRequest").(entity.DeviceEventRequest)
// 		if !ok {
// 			logger.Error("error retrieving deviceEventRequest from context")
// 			return echo.NewHTTPError(echo.ErrInternalServerError.Code, "internal server error, see logs.")
// 		}
// 		logger.Debug("income new event", zap.Any("msg", deviceEventRequest))

// 		err := s.deviceUseCase.SaveDeviceEventRequest(&deviceEventRequest)
// 		if err != nil {
// 			logger.Error("error SaveDeviceEventRequest", zap.Error(err))
// 			return echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to SaveDeviceEventRequest in handler.")
// 		}
// 		return ctx.NoContent(http.StatusOK)
// 	}
// }

// func (s *Server) GetDevices(logger *zap.Logger) echo.HandlerFunc {
// 	return func(ctx echo.Context) error {
// 		devices, err := s.deviceUseCase.GetDevices()
// 		if err != nil {
// 			logger.Error("error GetDevices", zap.Error(err))
// 			return echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to GetDevices in handler.")
// 		}
// 		logger.Debug("got GetDevices", zap.Any("msg", devices))

// 		if devices == nil {
// 			devices = &[]entity.DeviceInRepo{}
// 			return ctx.JSON(http.StatusNoContent, devices)
// 		}

// 		return ctx.JSON(http.StatusOK, devices)
// 	}
// }

// func (s *Server) GetDeviceEventsByID(logger *zap.Logger) echo.HandlerFunc {
// 	return func(ctx echo.Context) error {
// 		deviceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusBadRequest, "Invalid deviceID")
// 		}
// 		startTimeStr := ctx.QueryParam("start")
// 		endTimeStr := ctx.QueryParam("end")

// 		logger.Debug("got params", zap.Int64("deviceID", deviceID),
// 			zap.String("startTimeStr", startTimeStr),
// 			zap.String("endTimeStr", endTimeStr))

// 		// Prepare timestart to unixtime

// 		// Преобразование времени из RFC3339 в Unix time
// 		startTimeParsed, err := time.Parse(time.RFC3339, startTimeStr)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start time format")
// 		}
// 		startTimeUnix := startTimeParsed.Unix()

// 		endTimeParsed, err := time.Parse(time.RFC3339, endTimeStr)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusBadRequest, "Invalid end time format")
// 		}
// 		endTimeUnix := endTimeParsed.Unix()

// 		logger.Debug("parsed time", zap.Int64("startTimeUnix", startTimeUnix),
// 			zap.Int64("endTimeUnix", endTimeUnix))

// 		deviceEvents, err := s.deviceUseCase.GetDeviceEventsByID(deviceID, startTimeUnix, endTimeUnix)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error retrieving locations: %v", err))
// 		}
// 		if len(deviceEvents.Events) < 1 {
// 			logger.Debug("got empry locations", zap.Any("msg", deviceEvents))
// 			return ctx.JSON(http.StatusNoContent, deviceEvents)
// 		}
// 		logger.Debug("got locations", zap.Any("msg", deviceEvents))

// 		return ctx.JSON(http.StatusOK, deviceEvents)
// 	}
// }

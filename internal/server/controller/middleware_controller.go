package controller

import (
	"net/http"

	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func TokenAuthMiddleware(logger *zap.Logger, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logger.Error("Missing Authorization header")

				return c.JSON(http.StatusUnauthorized, Response{
					Data:  nil,
					Error: "Missing Authorization header",
				})
			}

			// Check if the token is in the format `Bearer <token>`
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token := authHeader[7:]

				// Validate the token (this is just a placeholder, you should implement your own validation logic)
				isValid, err := utils.IsValidToken(token, secret)

				if err != nil {
					return c.JSON(http.StatusInternalServerError, Response{
						Data:  nil,
						Error: "Error check auth token",
					})
				}

				if isValid {
					// If the token is valid, proceed with the next handler
					// userID, err := utils.GetUserIDFromToken(token, secret)
					if err != nil {
						logger.Error("Error get userID from token")

						return c.JSON(http.StatusUnauthorized, Response{
							Data:  nil,
							Error: "Error get userID from token",
						})
					}
					c.Set("user_id", "userID")

					return next(c)
				}
			}
			logger.Error("Invalid token")

			return c.JSON(http.StatusUnauthorized, Response{
				Data:  nil,
				Error: "Invalid token",
			})
		}
	}
}

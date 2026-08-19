package middleware

import (
	"net/http"
	"strings"

	"realworld-endpoints/internal/auth"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates JWT tokens from the Authorization header (supports "Token <jwt>" or "Bearer <jwt>")
func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"errors": echo.Map{"body": []string{"missing authorization header"}},
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || (parts[0] != "Token" && parts[0] != "Bearer") {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"errors": echo.Map{"body": []string{"invalid authorization header format"}},
				})
			}

			claims, err := auth.ParseToken(parts[1], jwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"errors": echo.Map{"body": []string{"invalid or expired token"}},
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)

			return next(c)
		}
	}
}

package middleware

import (
	"strings"

	"realworld-endpoints/internal/auth"

	"github.com/labstack/echo/v4"
)

// OptionalJWTMiddleware attempts to extract and validate JWT token if present in Authorization header.
// Does not reject requests if Authorization header is missing or invalid.
func OptionalJWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && (parts[0] == "Token" || parts[0] == "Bearer") {
					claims, err := auth.ParseToken(parts[1], jwtSecret)
					if err == nil {
						c.Set("user_id", claims.UserID)
						c.Set("user_email", claims.Email)
					}
				}
			}
			return next(c)
		}
	}
}

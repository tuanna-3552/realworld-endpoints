package routes

import (
	"net/http"

	"realworld-endpoints/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouterOptions struct {
	UserHandler    *handlers.UserHandler
	ArticleHandler *handlers.ArticleHandler
	TagHandler     *handlers.TagHandler
}

func SetupRoutes(e *echo.Echo, opts RouterOptions) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
	})

	// API Group
	api := e.Group("/api")

	// Users routes
	api.GET("/users", opts.UserHandler.GetUsers)

	// Articles routes
	api.GET("/articles", opts.ArticleHandler.GetArticles)

	// Tags routes
	api.GET("/tags", opts.TagHandler.GetTags)
}

package routes

import (
	"net/http"

	"realworld-endpoints/internal/handlers"
	customMiddleware "realworld-endpoints/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouterOptions struct {
	UserHandler    *handlers.UserHandler
	ArticleHandler *handlers.ArticleHandler
	TagHandler     *handlers.TagHandler
	ProfileHandler *handlers.ProfileHandler
	CommentHandler *handlers.CommentHandler
	JWTSecret      string
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

	// Authentication routes (Public)
	api.POST("/users", opts.UserHandler.Register)
	api.POST("/users/login", opts.UserHandler.Login)

	// Users routes (Public list)
	api.GET("/users", opts.UserHandler.GetUsers)

	// Articles routes
	api.GET("/articles", opts.ArticleHandler.GetArticles)
	api.GET("/articles/:slug", opts.ArticleHandler.GetArticleBySlug)

	// Comments routes
	api.GET("/articles/:slug/comments", opts.CommentHandler.GetComments)
	api.POST("/articles/:slug/comments", opts.CommentHandler.CreateComment)

	// Profiles routes
	api.GET("/profiles/:username", opts.ProfileHandler.GetProfile)

	// Tags routes
	api.GET("/tags", opts.TagHandler.GetTags)

	// Protected routes (Require JWT Authentication)
	authGroup := api.Group("", customMiddleware.JWTMiddleware(opts.JWTSecret))
	authGroup.GET("/user", opts.UserHandler.GetCurrentUser)
}

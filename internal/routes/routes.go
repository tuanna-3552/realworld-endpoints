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
	api.GET("/users", opts.UserHandler.GetUsers)

	// Tags routes (Public)
	api.GET("/tags", opts.TagHandler.GetTags)

	// Optional Auth Group (Parses token if provided, stays anonymous if omitted)
	optionalAuthGroup := api.Group("", customMiddleware.OptionalJWTMiddleware(opts.JWTSecret))
	optionalAuthGroup.GET("/articles", opts.ArticleHandler.GetArticles)
	optionalAuthGroup.GET("/articles/:slug", opts.ArticleHandler.GetArticleBySlug)
	optionalAuthGroup.GET("/articles/:slug/comments", opts.CommentHandler.GetComments)
	optionalAuthGroup.GET("/profiles/:username", opts.ProfileHandler.GetProfile)

	// Protected routes (Require JWT Authentication)
	authGroup := api.Group("", customMiddleware.JWTMiddleware(opts.JWTSecret))
	authGroup.GET("/user", opts.UserHandler.GetCurrentUser)

	// Step 5: Feed Articles
	authGroup.GET("/articles/feed", opts.ArticleHandler.GetFeed)

	// Step 6: Follow & Unfollow User
	authGroup.POST("/profiles/:username/follow", opts.ProfileHandler.FollowUser)
	authGroup.DELETE("/profiles/:username/follow", opts.ProfileHandler.UnfollowUser)

	// Step 6: Favorite & Unfavorite Article
	authGroup.POST("/articles/:slug/favorite", opts.ArticleHandler.FavoriteArticle)
	authGroup.DELETE("/articles/:slug/favorite", opts.ArticleHandler.UnfavoriteArticle)

	// Step 6: Create & Delete Comment (Protected + Ownership Check)
	authGroup.POST("/articles/:slug/comments", opts.CommentHandler.CreateComment)
	authGroup.DELETE("/articles/:slug/comments/:id", opts.CommentHandler.DeleteComment)
}

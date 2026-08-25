package main

import (
	"context"
	"log"

	"realworld-endpoints/internal/cache"
	"realworld-endpoints/internal/config"
	"realworld-endpoints/internal/db"
	"realworld-endpoints/internal/handlers"
	"realworld-endpoints/internal/repository"
	"realworld-endpoints/internal/routes"

	_ "realworld-endpoints/docs/swagger"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           RealWorld Conduit API
// @version         1.0
// @description     Backend API implementing the RealWorld spec built with Go, Echo Framework, GORM, and PostgreSQL.
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 JWT Token. Format: "Token {token}"
func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize Database (PostgreSQL via GORM & Seed Data)
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. Initialize Redis (graceful: warn and continue if unavailable)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	var cacheService cache.CacheService
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("[WARNING] Redis connection failed: %v. Caching disabled.", err)
		cacheService = nil
	} else {
		log.Println("[Redis] Connected successfully")
		cacheService = cache.NewRedisCacheService(redisClient)
	}

	// 4. Initialize Repositories
	userRepo := repository.NewUserRepository(database)
	articleRepo := repository.NewArticleRepository(database)
	tagRepo := repository.NewTagRepository(database)
	commentRepo := repository.NewCommentRepository(database)

	// 5. Initialize Handlers
	userHandler := handlers.NewUserHandler(userRepo, cfg.JWTSecret)
	articleHandler := handlers.NewArticleHandler(articleRepo, userRepo, cacheService)
	tagHandler := handlers.NewTagHandler(tagRepo, cacheService)
	profileHandler := handlers.NewProfileHandler(userRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo, articleRepo, userRepo)

	// 6. Initialize Echo Framework
	e := echo.New()

	// Swagger UI (development mode only)
	if cfg.IsDevelopment() {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		log.Println("[Swagger] UI enabled at /swagger/index.html (APP_ENV=development)")
	}

	// 7. Setup Routes
	routes.SetupRoutes(e, routes.RouterOptions{
		UserHandler:    userHandler,
		ArticleHandler: articleHandler,
		TagHandler:     tagHandler,
		ProfileHandler: profileHandler,
		CommentHandler: commentHandler,
		JWTSecret:      cfg.JWTSecret,
	})

	// 8. Start HTTP Server
	serverAddr := ":" + cfg.Port
	log.Printf("Starting Echo HTTP Server on %s...", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Server shutdown unexpectedly: %v", err)
	}
}

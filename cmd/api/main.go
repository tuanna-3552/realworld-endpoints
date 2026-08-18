package main

import (
	"log"

	"realworld-endpoints/internal/config"
	"realworld-endpoints/internal/db"
	"realworld-endpoints/internal/handlers"
	"realworld-endpoints/internal/repository"
	"realworld-endpoints/internal/routes"

	"github.com/labstack/echo/v4"
)

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

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(database)
	articleRepo := repository.NewArticleRepository(database)
	tagRepo := repository.NewTagRepository(database)
	commentRepo := repository.NewCommentRepository(database)

	// 4. Initialize Handlers
	userHandler := handlers.NewUserHandler(userRepo)
	articleHandler := handlers.NewArticleHandler(articleRepo)
	tagHandler := handlers.NewTagHandler(tagRepo)
	profileHandler := handlers.NewProfileHandler(userRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo, articleRepo, userRepo)

	// 5. Initialize Echo Framework
	e := echo.New()

	// 6. Setup Routes
	routes.SetupRoutes(e, routes.RouterOptions{
		UserHandler:    userHandler,
		ArticleHandler: articleHandler,
		TagHandler:     tagHandler,
		ProfileHandler: profileHandler,
		CommentHandler: commentHandler,
	})

	// 7. Start HTTP Server
	serverAddr := ":" + cfg.Port
	log.Printf("Starting Echo HTTP Server on %s...", serverAddr)
	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Server shutdown unexpectedly: %v", err)
	}
}

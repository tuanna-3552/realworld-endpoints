package db

import (
	"fmt"
	"log"
	"time"

	"realworld-endpoints/internal/config"
	"realworld-endpoints/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DSN()
	var db *gorm.DB
	var err error

	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := 1; i <= maxRetries; i++ {
		log.Printf("Connecting to PostgreSQL (Attempt %d/%d)...", i, maxRetries)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			// Verify raw connection works
			sqlDB, pingErr := db.DB()
			if pingErr == nil {
				if pingErr = sqlDB.Ping(); pingErr == nil {
					log.Println("Connected to PostgreSQL successfully.")
					break
				}
			}
			err = pingErr
		}

		log.Printf("PostgreSQL connection failed: %v. Retrying in %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database after %d attempts: %w", maxRetries, err)
	}

	log.Println("Running AutoMigrate...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Tag{},
		&models.Article{},
		&models.Comment{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}

	log.Println("Database AutoMigrate completed successfully.")

	// Automatically seed initial data if DB is empty
	if seedErr := Seed(db); seedErr != nil {
		log.Printf("Warning: Database seeding failed: %v", seedErr)
	}

	return db, nil
}

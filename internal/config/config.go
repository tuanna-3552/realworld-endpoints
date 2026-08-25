package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string

	// Application Environment
	AppEnv string

	// Redis Configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() (*Config, error) {
	// Flexible .env loading: .env.local overrides .env
	// Load .env.local first (personal overrides, gitignored)
	_ = godotenv.Load(".env.local")
	// Then load .env (shared defaults)
	_ = godotenv.Load(".env")

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "54333"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "1c0b1c9e-59c4-4c26-896f-d4a5795a2c9a"),
		DBName:     getEnv("DB_NAME", "realworld"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "secret-jwt-key-change-in-production"),
		AppEnv:     getEnv("APP_ENV", "development"),

		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	log.Printf("[Config] Loaded successfully: Port=%s, DBHost=%s, DBName=%s, RedisHost=%s:%s",
		cfg.Port, cfg.DBHost, cfg.DBName, cfg.RedisHost, cfg.RedisPort)

	return cfg, nil
}

// Validate checks that all required configuration fields are set
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("config validation failed: PORT is required")
	}
	if c.DBHost == "" {
		return fmt.Errorf("config validation failed: DB_HOST is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("config validation failed: DB_NAME is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("config validation failed: JWT_SECRET is required")
	}
	return nil
}

// DSN returns the PostgreSQL connection string
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// RedisAddr returns the Redis server address in "host:port" format
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

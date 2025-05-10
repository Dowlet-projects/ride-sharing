// config/config.go
// Configuration loading for the application

package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	JWTSecret  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"), // Can be empty
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	// Validate required fields (DB_PASSWORD can be empty)
	if cfg.DBUser == "" || cfg.DBName == "" || cfg.DBHost == "" || cfg.JWTSecret == "" {
		return nil, fmt.Errorf("missing required environment variables: DB_USER=%q, DB_NAME=%q, DB_HOST=%q, JWT_SECRET=%q",
			cfg.DBUser, cfg.DBName, cfg.DBHost, cfg.JWTSecret)
	}

	return cfg, nil
}
// Package config loads and validates application configuration from environment variables.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for the application.
type Config struct {
	Port        string
	DatabaseURL string
	LogLevel    string
	Env         string
}

// Load reads configuration from environment variables.
// If a .env file is present it is loaded first (development convenience).
func Load() (*Config, error) {
	// Ignore error — .env is optional (not present in production).
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Env:         getEnv("ENV", "dev"),
	}

	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL is required but not set")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

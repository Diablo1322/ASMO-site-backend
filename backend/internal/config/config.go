package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	FrontendURL string
	LogLevel    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/ASMO-site-backenddb"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		LogLevel:    getEnv("LOG_LEVEL", "INFO"),
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
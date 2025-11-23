package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	LogLevel    string
	Environment string
}

func Load() *Config {
	environment := getEnv("ENVIRONMENT", "development")

	return &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/asmo_db?sslmode=disable"),
		LogLevel:    getEnv("LOG_LEVEL", getDefaultLogLevel(environment)),
		Environment: environment,
	}
}

func getDefaultLogLevel(environment string) string {
	if environment == "production" {
		return "INFO"
	}
	return "DEBUG"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
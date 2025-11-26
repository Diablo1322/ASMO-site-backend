package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	DatabaseURL    string
	LogLevel       string
	Environment    string
	AllowedOrigins string
}

func Load() *Config {
	environment := getEnv("ENVIRONMENT", "development")

	return &Config{
		Port:           getEnv("PORT", "3000"),
		DatabaseURL:    getDatabaseURL(environment),
		LogLevel:       getEnv("LOG_LEVEL", getDefaultLogLevel(environment)),
		Environment:    environment,
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", getAllowedOrigins(environment)),
	}
}

func getDatabaseURL(environment string) string {
	if environment == "production" {
		host := getEnv("DB_HOST", "postgres")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "asmo_prod_user")
		password := getEnv("DB_PASSWORD", "")
		dbname := getEnv("DB_NAME", "asmo_production")
		sslMode := getEnv("DB_SSL_MODE", "require")

		return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslMode
	}

	return getEnv("DATABASE_URL", "postgres://user:password@postgres:5432/asmo_db?sslmode=disable")
}

func getAllowedOrigins(environment string) string {
	if environment == "production" {
		return getEnv("ALLOWED_ORIGINS", "https://production-domain.com")
	}
	return getEnv("ALLOWED_ORIGINS", "http://localhost:3001,http://127.0.0.1:3001")
}

func getDefaultLogLevel(environment string) string {
	if environment == "production" {
		return "INFO"
	}
	return "DEBUG"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}
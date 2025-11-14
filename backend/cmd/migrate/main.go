package main

import (
	"log"
	"os"
	"strconv"

	"ASMO-site-backend/internal/config"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	cfg := config.Load()

	// Создание миграций
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal("Failed to create migrator:", err)
	}

	// Обработка аргументов командной строки
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down|version|force VERSION]")
	}

	command := os.Args[1]

	switch command {
	case "up":
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to migrate up:", err)
		}
		log.Println("Migration up completed successfully")

	case "down":
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to migrate down:", err)
		}
		log.Println("Migration down completed successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get version:", err)
		}
		log.Printf("Current version: %d, dirty: %t", version, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Force command requires a version number")
		}
		versionStr := os.Args[2]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			log.Fatal("Invalid version number:", versionStr)
		}
		err = m.Force(version)
		if err != nil {
			log.Fatal("Failed to force version:", err)
		}
		log.Printf("Forced version to: %d", version)

	default:
		log.Fatal("Unknown command:", command)
	}
}
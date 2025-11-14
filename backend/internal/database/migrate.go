package database

import (
	"log"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

// RunMigrations автоматически применяет миграции при запуске приложения
func RunMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return err
	}

	// Применяем все pending миграции
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	if err == migrate.ErrNilVersion {
		log.Println("No migrations applied - database is empty")
	} else {
		log.Printf("Migrations up to date. Version: %d, Dirty: %t", version, dirty)
	}

	return nil
}
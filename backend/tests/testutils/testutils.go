package testutils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

// SetupTestDB создает тестовую базу данных и применяет миграции
func SetupTestDB() (*sql.DB, error) {
	connStr := "postgres://test:test@localhost:5433/asmo_test_db?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	// Применяем миграции
	m, err := migrate.New(
		"file://../../migrations",
		connStr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Test database setup completed successfully")
	return db, nil
}

// CleanupTestDB очищает тестовую базу данных
func CleanupTestDB(db *sql.DB) error {
	// Откатываем все миграции
	m, err := migrate.New(
		"file://../../migrations",
		"postgres://test:test@localhost:5433/asmo_test_db?sslmode=disable",
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator for cleanup: %w", err)
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return db.Close()
}

package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

// getMigrationsPath возвращает абсолютный путь к папке миграций
func getMigrationsPath() (string, error) {
	// Получаем текущую рабочую директорию
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Проверяем разные возможные расположения
	possiblePaths := []string{
		filepath.Join(wd, "migrations"),           // из корня backend
		filepath.Join(wd, "..", "migrations"),     // из tests/integration
		filepath.Join(wd, "..", "..", "migrations"), // из tests
		filepath.Join(wd, "backend", "migrations"), // из корня проекта
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("Found migrations at: %s", path)
			return path, nil
		}
	}

	return "", fmt.Errorf("migrations directory not found. Checked paths: %v", possiblePaths)
}

// waitForDB ждет пока база данных станет доступной
func waitForDB(connStr string, timeout time.Duration) error {
	start := time.Now()
	for {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return err
		}
		defer db.Close()

		err = db.Ping()
		if err == nil {
			log.Println("Database connection established")
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for database after %v", timeout)
		}

		log.Printf("Waiting for database... (%v)", err)
		time.Sleep(2 * time.Second)
	}
}

// SetupTestDB создает тестовую базу данных и применяет миграции
func SetupTestDB() (*sql.DB, error) {
	connStr := "postgres://test:test@localhost:5433/testdb?sslmode=disable"

	// Ждем пока база данных станет доступной
	log.Println("Waiting for test database...")
	err := waitForDB(connStr, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	// Получаем путь к миграциям
	migrationsPath, err := getMigrationsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find migrations: %w", err)
	}

	// Конвертируем путь в формат file:// для migrate
	migrationsURL := "file://" + filepath.ToSlash(migrationsPath)

	log.Printf("Applying migrations from: %s", migrationsURL)

	// Применяем миграции
	m, err := migrate.New(migrationsURL, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Migration version: %d, dirty: %t", version, dirty)
	}

	log.Println("Test database setup completed successfully")
	return db, nil
}

// CleanupTestDB очищает тестовую базу данных
func CleanupTestDB(db *sql.DB) error {
	connStr := "postgres://test:test@localhost:5433/testdb?sslmode=disable"

	// Получаем путь к миграциям
	migrationsPath, err := getMigrationsPath()
	if err != nil {
		return fmt.Errorf("failed to find migrations: %w", err)
	}

	migrationsURL := "file://" + filepath.ToSlash(migrationsPath)

	// Откатываем все миграции
	m, err := migrate.New(migrationsURL, connStr)
	if err != nil {
		return fmt.Errorf("failed to create migrator for cleanup: %w", err)
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Println("Test database cleanup completed")
	return db.Close()
}
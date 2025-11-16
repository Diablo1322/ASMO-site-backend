package testutils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func SetupTestDB() (*sql.DB, error) {
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://test:test@localhost:5433/testdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	// Create test tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create test tables: %w", err)
	}

	return db, nil
}

func CleanupTestDB(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS items")
	if err != nil {
		return fmt.Errorf("failed to clean test database: %w", err)
	}
	return db.Close()
}
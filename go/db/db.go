package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(path string) (*sql.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("database path cannot be empty")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("creating database directory: %w", err)
	}

	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// SQLite admite un único escritor simultáneo.
	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	if _, err := database.Exec("PRAGMA journal_mode = WAL"); err != nil {
		database.Close()
		return nil, fmt.Errorf("configuring WAL mode: %w", err)
	}

	if _, err := database.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		database.Close()
		return nil, fmt.Errorf("configuring busy timeout: %w", err)
	}

	if err := CreateSchema(database); err != nil {
		database.Close()
		return nil, fmt.Errorf("creating database schema: %w", err)
	}

	return database, nil
}

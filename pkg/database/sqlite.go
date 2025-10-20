package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Connect opens a connection to the SQLite database
func Connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Close safely closes the database connection
func Close(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".note")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("could not create directory %s: %w", dbDir, err)
	}

	dbPath := filepath.Join(dbDir, "notes.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("could not connect to database at %s: %w", dbPath, err)
	}

	DB = db

	if err := runMigrations(db); err != nil {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	if err := ensureHomeProject(db); err != nil {
		return fmt.Errorf("could not ensure home project exists: %w", err)
	}

	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

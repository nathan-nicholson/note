package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	schemas := []string{
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			is_important BOOLEAN NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			is_complete BOOLEAN NOT NULL DEFAULT 0,
			due_date DATE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS note_tags (
			note_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (note_id, tag_id)
		)`,
		`CREATE TABLE IF NOT EXISTS todo_tags (
			todo_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (todo_id, tag_id)
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			first_activated_at TIMESTAMP,
			last_activity_at TIMESTAMP,
			closed_at TIMESTAMP,
			is_closed BOOLEAN DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS project_tags (
			project_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (project_id, tag_id)
		)`,
		`CREATE TABLE IF NOT EXISTS active_project (
			project_id INTEGER NOT NULL UNIQUE,
			activated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,
	}

	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			t.Fatalf("Failed to create test schema: %v", err)
		}
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

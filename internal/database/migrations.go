package database

import (
	"database/sql"
	"fmt"
)

func runMigrations(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_important BOOLEAN NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		is_complete BOOLEAN NOT NULL DEFAULT 0,
		due_date DATE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS note_tags (
		note_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
		PRIMARY KEY (note_id, tag_id)
	);

	CREATE TABLE IF NOT EXISTS todo_tags (
		todo_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
		PRIMARY KEY (todo_id, tag_id)
	);

	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		first_activated_at TIMESTAMP,
		last_activity_at TIMESTAMP,
		closed_at TIMESTAMP,
		is_closed BOOLEAN DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS project_tags (
		project_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
		PRIMARY KEY (project_id, tag_id)
	);

	CREATE TABLE IF NOT EXISTS active_project (
		project_id INTEGER NOT NULL UNIQUE,
		activated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func ensureHomeProject(db *sql.DB) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE name = 'home')").Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		result, err := tx.Exec(`
			INSERT INTO projects (name, created_at, first_activated_at, is_closed)
			VALUES ('home', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
		`)
		if err != nil {
			return err
		}

		projectID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO active_project (project_id, activated_at)
			VALUES (?, CURRENT_TIMESTAMP)
		`, projectID)
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

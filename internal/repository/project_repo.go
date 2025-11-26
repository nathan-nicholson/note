package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nathan-nicholson/note/internal/models"
)

func CreateProject(db *sql.DB, name string, tags []string) (*models.Project, error) {
	if err := models.ValidateProjectName(name); err != nil {
		return nil, err
	}

	now := time.Now()
	result, err := db.Exec(`
		INSERT INTO projects (name, created_at, is_closed)
		VALUES (?, ?, 0)
	`, name, now)
	if err != nil {
		return nil, err
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := AddTagsToProject(db, int(projectID), tags); err != nil {
		return nil, err
	}

	return GetProjectByName(db, name)
}

func GetProjectByName(db *sql.DB, name string) (*models.Project, error) {
	var project models.Project
	err := db.QueryRow(`
		SELECT id, name, created_at, first_activated_at, last_activity_at, closed_at, is_closed
		FROM projects
		WHERE name = ?
	`, name).Scan(&project.ID, &project.Name, &project.CreatedAt, &project.FirstActivatedAt,
		&project.LastActivityAt, &project.ClosedAt, &project.IsClosed)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Project '%s' not found. Create it with: note project create %s", name, name)
		}
		return nil, err
	}

	tags, err := GetTagsForProject(db, project.ID)
	if err != nil {
		return nil, err
	}
	project.Tags = tags

	return &project, nil
}

func GetProjectByID(db *sql.DB, id int) (*models.Project, error) {
	var project models.Project
	err := db.QueryRow(`
		SELECT id, name, created_at, first_activated_at, last_activity_at, closed_at, is_closed
		FROM projects
		WHERE id = ?
	`, id).Scan(&project.ID, &project.Name, &project.CreatedAt, &project.FirstActivatedAt,
		&project.LastActivityAt, &project.ClosedAt, &project.IsClosed)

	if err != nil {
		return nil, err
	}

	tags, err := GetTagsForProject(db, project.ID)
	if err != nil {
		return nil, err
	}
	project.Tags = tags

	return &project, nil
}

func GetActiveProject(db *sql.DB) (*models.Project, error) {
	var projectID int
	err := db.QueryRow("SELECT project_id FROM active_project").Scan(&projectID)
	if err != nil {
		return nil, err
	}

	return GetProjectByID(db, projectID)
}

func SetActiveProject(db *sql.DB, projectID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	_, err = tx.Exec("DELETE FROM active_project")
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO active_project (project_id, activated_at)
		VALUES (?, ?)
	`, projectID, now)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE projects
		SET first_activated_at = COALESCE(first_activated_at, ?)
		WHERE id = ?
	`, now, projectID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func ListProjects(db *sql.DB, includeAll bool) ([]models.Project, error) {
	query := `
		SELECT id, name, created_at, first_activated_at, last_activity_at, closed_at, is_closed
		FROM projects
	`

	if !includeAll {
		query += " WHERE is_closed = 0"
	}

	query += " ORDER BY name"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.CreatedAt, &project.FirstActivatedAt,
			&project.LastActivityAt, &project.ClosedAt, &project.IsClosed); err != nil {
			return nil, err
		}

		tags, err := GetTagsForProject(db, project.ID)
		if err != nil {
			return nil, err
		}
		project.Tags = tags

		projects = append(projects, project)
	}

	return projects, rows.Err()
}

func CloseProject(db *sql.DB, projectID int) error {
	now := time.Now()
	_, err := db.Exec(`
		UPDATE projects
		SET is_closed = 1, closed_at = ?
		WHERE id = ?
	`, now, projectID)
	return err
}

func ReopenProject(db *sql.DB, projectID int) error {
	_, err := db.Exec(`
		UPDATE projects
		SET is_closed = 0, closed_at = NULL
		WHERE id = ?
	`, projectID)
	return err
}

func DeleteProject(db *sql.DB, name string) error {
	if name == "home" {
		return fmt.Errorf("Cannot delete the 'home' project. It is the default project and must always exist.")
	}

	result, err := db.Exec("DELETE FROM projects WHERE name = ?", name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Project '%s' not found", name)
	}

	return nil
}

func UpdateProjectTags(db *sql.DB, projectID int, tags []string) error {
	return ReplaceProjectTags(db, projectID, tags)
}

func GetIncompleteTodosForProject(db *sql.DB, projectName string) ([]models.Todo, error) {
	rows, err := db.Query(`
		SELECT DISTINCT t.id, t.content, t.is_complete, t.due_date, t.created_at, t.updated_at, t.completed_at
		FROM todos t
		JOIN todo_tags tt ON t.id = tt.todo_id
		JOIN tags tg ON tt.tag_id = tg.id
		WHERE tg.name = ? AND t.is_complete = 0
		ORDER BY t.due_date IS NULL, t.due_date, t.created_at
	`, projectName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Content, &todo.IsComplete, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt, &todo.CompletedAt); err != nil {
			return nil, err
		}

		tags, err := GetTagsForTodo(db, todo.ID)
		if err != nil {
			return nil, err
		}
		todo.Tags = tags

		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func GetCompleteTodosForProject(db *sql.DB, projectName string) ([]models.Todo, error) {
	rows, err := db.Query(`
		SELECT DISTINCT t.id, t.content, t.is_complete, t.due_date, t.created_at, t.updated_at, t.completed_at
		FROM todos t
		JOIN todo_tags tt ON t.id = tt.todo_id
		JOIN tags tg ON tt.tag_id = tg.id
		WHERE tg.name = ? AND t.is_complete = 1
		ORDER BY t.completed_at DESC
	`, projectName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Content, &todo.IsComplete, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt, &todo.CompletedAt); err != nil {
			return nil, err
		}

		tags, err := GetTagsForTodo(db, todo.ID)
		if err != nil {
			return nil, err
		}
		todo.Tags = tags

		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func UpdateProjectLastActivity(db *sql.DB, projectID int) error {
	now := time.Now()
	_, err := db.Exec(`
		UPDATE projects
		SET last_activity_at = ?
		WHERE id = ?
	`, now, projectID)
	return err
}

func CountOpenProjects(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM projects WHERE is_closed = 0").Scan(&count)
	return count, err
}

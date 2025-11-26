package activity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/nathan-nicholson/note/internal/models"
	"github.com/nathan-nicholson/note/internal/repository"
)

func LogTodoCreated(db *sql.DB, todo *models.Todo) error {
	content := fmt.Sprintf("Created todo: %s", todo.Content)

	if todo.DueDate.Valid {
		content += fmt.Sprintf(" (due: %s)", todo.DueDate.Time.Format("2006-01-02"))
	}

	tags := append([]string{"todo", "create"}, todo.Tags...)
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogTodoUpdated(db *sql.DB, todo *models.Todo, changes []string) error {
	if len(changes) == 0 {
		return nil
	}

	content := fmt.Sprintf("Updated todo: %s", strings.Join(changes, ", "))

	tags := []string{"todo", "update"}
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogTodoCompleted(db *sql.DB, todo *models.Todo) error {
	content := fmt.Sprintf("Completed todo: %s", todo.Content)

	if todo.DueDate.Valid {
		content += fmt.Sprintf(" (due: %s)", todo.DueDate.Time.Format("2006-01-02"))
	}

	tags := append([]string{"todo", "complete"}, todo.Tags...)
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogTodoDeleted(db *sql.DB, todo *models.Todo) error {
	content := fmt.Sprintf("Deleted todo: %s", todo.Content)

	if todo.DueDate.Valid {
		content += fmt.Sprintf(" (due: %s)", todo.DueDate.Time.Format("2006-01-02"))
	}

	tags := append([]string{"todo", "delete"}, todo.Tags...)
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogProjectCreated(db *sql.DB, project *models.Project) error {
	content := fmt.Sprintf("Created project: %s", project.Name)

	tags := append([]string{"project", "create"}, project.Tags...)
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogProjectActivated(db *sql.DB, projectName string) error {
	content := fmt.Sprintf("Activated project: %s", projectName)
	tags := []string{"project", "activate"}
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogProjectDeactivated(db *sql.DB, projectName string) error {
	content := fmt.Sprintf("Deactivated project: %s", projectName)
	tags := []string{"project", "deactivate"}
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogProjectUpdated(db *sql.DB, projectName string, changes string) error {
	content := fmt.Sprintf("Updated project: %s - %s", projectName, changes)
	tags := []string{"project", "update"}
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogProjectClosed(db *sql.DB, projectName string) error {
	content := fmt.Sprintf("Closed project: %s", projectName)
	tags := []string{"project", "close"}
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogProjectReopened(db *sql.DB, projectName string) error {
	content := fmt.Sprintf("Reopened project: %s", projectName)
	tags := []string{"project", "reopen"}
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

func LogProjectDeleted(db *sql.DB, project *models.Project) error {
	content := fmt.Sprintf("Deleted project: %s", project.Name)
	tags := append([]string{"project", "delete"}, project.Tags...)
	_, err := repository.CreateNote(db, content, tags, false)
	return err
}

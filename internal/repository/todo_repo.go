package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/nathan-nicholson/note/internal/models"
)

func CreateTodo(db *sql.DB, content string, tags []string, dueDate *time.Time) (*models.Todo, error) {
	var dueDateSQL interface{}
	if dueDate != nil {
		dueDateSQL = dueDate.Format("2006-01-02")
	}

	result, err := db.Exec(`
		INSERT INTO todos (content, due_date, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, content, dueDateSQL)
	if err != nil {
		return nil, err
	}

	todoID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := AddTagsToTodo(db, int(todoID), tags); err != nil {
		return nil, err
	}

	return GetTodoByID(db, int(todoID))
}

func GetTodoByID(db *sql.DB, id int) (*models.Todo, error) {
	var todo models.Todo
	err := db.QueryRow(`
		SELECT id, content, is_complete, due_date, created_at, updated_at, completed_at
		FROM todos
		WHERE id = ?
	`, id).Scan(&todo.ID, &todo.Content, &todo.IsComplete, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt, &todo.CompletedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Todo #%d not found", id)
		}
		return nil, err
	}

	tags, err := GetTagsForTodo(db, id)
	if err != nil {
		return nil, err
	}
	todo.Tags = tags

	return &todo, nil
}

type TodoListOptions struct {
	Complete   bool
	Incomplete bool
	Tags       []string
	Overdue    bool
}

func ListTodos(db *sql.DB, opts TodoListOptions) ([]models.Todo, error) {
	query := `
		SELECT DISTINCT t.id, t.content, t.is_complete, t.due_date, t.created_at, t.updated_at, t.completed_at
		FROM todos t
	`

	var conditions []string
	var args []interface{}

	if len(opts.Tags) > 0 {
		query += `
			JOIN todo_tags tt ON t.id = tt.todo_id
			JOIN tags tg ON tt.tag_id = tg.id
		`
		placeholders := make([]string, len(opts.Tags))
		for i, tag := range opts.Tags {
			placeholders[i] = "?"
			args = append(args, tag)
		}
		conditions = append(conditions, fmt.Sprintf("tg.name IN (%s)", strings.Join(placeholders, ",")))
	}

	if opts.Complete {
		conditions = append(conditions, "t.is_complete = 1")
	}

	if opts.Incomplete {
		conditions = append(conditions, "t.is_complete = 0")
	}

	if opts.Overdue {
		conditions = append(conditions, "t.due_date IS NOT NULL AND DATE(t.due_date) < DATE('now') AND t.is_complete = 0")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if len(opts.Tags) > 0 {
		query += fmt.Sprintf(" GROUP BY t.id HAVING COUNT(DISTINCT tg.name) = %d", len(opts.Tags))
	}

	query += " ORDER BY t.due_date IS NULL, t.due_date, t.created_at"

	rows, err := db.Query(query, args...)
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

func CompleteTodo(db *sql.DB, id int) error {
	_, err := db.Exec(`
		UPDATE todos
		SET is_complete = 1, completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)
	return err
}

func UncompleteTodo(db *sql.DB, id int) error {
	_, err := db.Exec(`
		UPDATE todos
		SET is_complete = 0, completed_at = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)
	return err
}

func UpdateTodo(db *sql.DB, id int, content *string, tags []string, dueDate *time.Time, clearDueDate bool) error {
	if content != nil {
		_, err := db.Exec(`
			UPDATE todos
			SET content = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, *content, id)
		if err != nil {
			return err
		}
	}

	if clearDueDate {
		_, err := db.Exec(`
			UPDATE todos
			SET due_date = NULL, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, id)
		if err != nil {
			return err
		}
	} else if dueDate != nil {
		_, err := db.Exec(`
			UPDATE todos
			SET due_date = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, dueDate.Format("2006-01-02"), id)
		if err != nil {
			return err
		}
	}

	if len(tags) > 0 {
		if err := ReplaceTodoTags(db, id, tags); err != nil {
			return err
		}
		_, err := db.Exec("UPDATE todos SET updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteTodo(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Todo #%d not found", id)
	}

	return nil
}

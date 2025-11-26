package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/nathan-nicholson/note/internal/models"
)

type TodoGroup struct {
	Title string
	Todos []models.Todo
}

func GroupTodos(todos []models.Todo) []TodoGroup {
	today := time.Now().Truncate(24 * time.Hour)

	var overdue []models.Todo
	var todayTodos []models.Todo
	var upcoming []models.Todo
	var noDueDate []models.Todo

	for _, todo := range todos {
		if !todo.DueDate.Valid {
			noDueDate = append(noDueDate, todo)
		} else {
			dueDate := todo.DueDate.Time.Truncate(24 * time.Hour)

			if dueDate.Before(today) && !todo.IsComplete {
				overdue = append(overdue, todo)
			} else if dueDate.Equal(today) {
				todayTodos = append(todayTodos, todo)
			} else if dueDate.After(today) {
				upcoming = append(upcoming, todo)
			}
		}
	}

	var groups []TodoGroup

	if len(overdue) > 0 {
		groups = append(groups, TodoGroup{Title: "OVERDUE", Todos: overdue})
	}

	if len(todayTodos) > 0 {
		groups = append(groups, TodoGroup{
			Title: fmt.Sprintf("TODAY (%s)", time.Now().Format("2006-01-02")),
			Todos: todayTodos,
		})
	}

	if len(upcoming) > 0 {
		groups = append(groups, TodoGroup{Title: "UPCOMING", Todos: upcoming})
	}

	if len(noDueDate) > 0 {
		groups = append(groups, TodoGroup{Title: "NO DUE DATE", Todos: noDueDate})
	}

	return groups
}

func FormatTodoList(todos []models.Todo) string {
	if len(todos) == 0 {
		return ""
	}

	groups := GroupTodos(todos)

	var output strings.Builder

	for i, group := range groups {
		if i > 0 {
			output.WriteString("\n")
		}

		output.WriteString(group.Title + "\n")

		for _, todo := range group.Todos {
			output.WriteString("  ")

			if todo.IsComplete {
				output.WriteString("[X]")
			} else {
				output.WriteString("[ ]")
			}

			output.WriteString(fmt.Sprintf(" [#%d] ", todo.ID))

			if todo.DueDate.Valid && group.Title != "NO DUE DATE" {
				dueDate := todo.DueDate.Time.Format("2006-01-02")
				if group.Title == "UPCOMING" {
					output.WriteString(dueDate + "  ")
				} else if strings.HasPrefix(group.Title, "TODAY") {
					output.WriteString("(due today) ")
				} else if group.Title == "OVERDUE" {
					output.WriteString(dueDate + "  ")
				}
			}

			output.WriteString(todo.Content)

			if len(todo.Tags) > 0 {
				output.WriteString(" ")
				for _, tag := range todo.Tags {
					output.WriteString("#" + tag + " ")
				}
			}

			output.WriteString("\n")
		}
	}

	return strings.TrimSpace(output.String())
}

func FormatTodo(todo *models.Todo) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("ID: %d\n", todo.ID))
	output.WriteString(fmt.Sprintf("Created: %s\n", todo.CreatedAt.Format("2006-01-02 03:04 PM")))

	if todo.DueDate.Valid {
		output.WriteString(fmt.Sprintf("Due: %s\n", todo.DueDate.Time.Format("2006-01-02")))
	}

	output.WriteString(fmt.Sprintf("Status: %s\n", map[bool]string{true: "Complete", false: "Incomplete"}[todo.IsComplete]))

	if todo.CompletedAt.Valid {
		output.WriteString(fmt.Sprintf("Completed: %s\n", todo.CompletedAt.Time.Format("2006-01-02 03:04 PM")))
	}

	if len(todo.Tags) > 0 {
		output.WriteString("Tags: ")
		for i, tag := range todo.Tags {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString("#" + tag)
		}
		output.WriteString("\n")
	}

	output.WriteString("\n")
	output.WriteString(todo.Content)

	return output.String()
}

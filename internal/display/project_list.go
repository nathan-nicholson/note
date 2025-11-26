package display

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/nathan-nicholson/note/internal/models"
	"github.com/nathan-nicholson/note/internal/repository"
)

func FormatProjectList(db *sql.DB, projects []models.Project, includeAll bool) (string, error) {
	if len(projects) == 0 {
		return "", nil
	}

	activeProject, err := repository.GetActiveProject(db)
	if err != nil {
		return "", err
	}

	var active []models.Project
	var open []models.Project
	var closed []models.Project

	for _, project := range projects {
		if project.ID == activeProject.ID {
			active = append(active, project)
		} else if project.IsClosed {
			closed = append(closed, project)
		} else {
			open = append(open, project)
		}
	}

	var output strings.Builder

	if len(active) > 0 {
		output.WriteString("ACTIVE\n")
		for _, p := range active {
			output.WriteString("  * " + p.Name)
			if len(p.Tags) > 0 {
				output.WriteString(" ")
				for _, tag := range p.Tags {
					output.WriteString("#" + tag + " ")
				}
			}
			output.WriteString("\n")
		}
	}

	if len(open) > 0 {
		if output.Len() > 0 {
			output.WriteString("\n")
		}
		output.WriteString("OPEN\n")
		for _, p := range open {
			output.WriteString("  " + p.Name)
			if len(p.Tags) > 0 {
				output.WriteString(" ")
				for _, tag := range p.Tags {
					output.WriteString("#" + tag + " ")
				}
			}
			output.WriteString("\n")
		}
	}

	if includeAll && len(closed) > 0 {
		if output.Len() > 0 {
			output.WriteString("\n")
		}
		output.WriteString("CLOSED\n")
		for _, p := range closed {
			output.WriteString("  " + p.Name)
			if len(p.Tags) > 0 {
				output.WriteString(" ")
				for _, tag := range p.Tags {
					output.WriteString("#" + tag + " ")
				}
			}
			output.WriteString("\n")
		}
	}

	return strings.TrimSpace(output.String()), nil
}

func FormatProjectStatus(db *sql.DB, project *models.Project, showAll bool) (string, error) {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("Project: %s\n", project.Name))

	if len(project.Tags) > 0 {
		output.WriteString("Tags: ")
		for i, tag := range project.Tags {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString("#" + tag)
		}
		output.WriteString("\n")
	}

	if project.IsClosed {
		output.WriteString("Status: Closed\n")
	} else {
		output.WriteString("Status: Open\n")
	}

	output.WriteString(fmt.Sprintf("Created: %s\n", project.CreatedAt.Format("2006-01-02")))

	if project.FirstActivatedAt.Valid {
		output.WriteString(fmt.Sprintf("First Activated: %s\n", project.FirstActivatedAt.Time.Format("2006-01-02")))
	}

	if project.LastActivityAt.Valid {
		output.WriteString(fmt.Sprintf("Last Activity: %s\n", project.LastActivityAt.Time.Format("2006-01-02")))
	}

	if project.ClosedAt.Valid {
		output.WriteString(fmt.Sprintf("Closed: %s\n", project.ClosedAt.Time.Format("2006-01-02")))
	}

	incompleteTodos, err := repository.GetIncompleteTodosForProject(db, project.Name)
	if err != nil {
		return "", err
	}

	completeTodos, err := repository.GetCompleteTodosForProject(db, project.Name)
	if err != nil {
		return "", err
	}

	totalTodos := len(incompleteTodos) + len(completeTodos)
	completedCount := len(completeTodos)

	var percentage int
	if totalTodos > 0 {
		percentage = (completedCount * 100) / totalTodos
	}

	output.WriteString(fmt.Sprintf("Tasks: %d/%d complete (%d%%)\n", completedCount, totalTodos, percentage))

	if len(incompleteTodos) > 0 {
		output.WriteString("\nIncomplete Tasks:\n")
		for _, todo := range incompleteTodos {
			output.WriteString(fmt.Sprintf("  [ ] [#%d] ", todo.ID))
			if todo.DueDate.Valid {
				output.WriteString(fmt.Sprintf("%s  ", todo.DueDate.Time.Format("2006-01-02")))
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

	if showAll && len(completeTodos) > 0 {
		output.WriteString("\nCompleted Tasks:\n")
		for _, todo := range completeTodos {
			output.WriteString(fmt.Sprintf("  [X] [#%d] ", todo.ID))
			if todo.DueDate.Valid {
				output.WriteString(fmt.Sprintf("(was due: %s) ", todo.DueDate.Time.Format("2006-01-02")))
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

	return strings.TrimSpace(output.String()), nil
}

func FormatProject(project *models.Project) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("Project: %s\n", project.Name))

	if len(project.Tags) > 0 {
		output.WriteString("Tags: ")
		for i, tag := range project.Tags {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString("#" + tag)
		}
		output.WriteString("\n")
	}

	if project.IsClosed {
		output.WriteString("Status: Closed\n")
	} else {
		output.WriteString("Status: Open\n")
	}

	output.WriteString(fmt.Sprintf("Created: %s\n", project.CreatedAt.Format("2006-01-02")))

	if project.FirstActivatedAt.Valid {
		output.WriteString(fmt.Sprintf("First Activated: %s\n", project.FirstActivatedAt.Time.Format("2006-01-02")))
	}

	if project.LastActivityAt.Valid {
		output.WriteString(fmt.Sprintf("Last Activity: %s\n", project.LastActivityAt.Time.Format("2006-01-02")))
	}

	if project.ClosedAt.Valid {
		output.WriteString(fmt.Sprintf("Closed: %s\n", project.ClosedAt.Time.Format("2006-01-02")))
	}

	return strings.TrimSpace(output.String())
}

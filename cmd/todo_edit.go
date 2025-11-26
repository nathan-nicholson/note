package cmd

import (
	"strconv"
	"strings"
	"time"

	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/dateparse"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var (
	todoEditContent string
	todoEditTags    []string
	todoEditDue     string
)

var todoEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a todo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		oldTodo, err := repository.GetTodoByID(database.DB, id)
		if err != nil {
			return err
		}

		var content *string
		var dueDate *time.Time
		clearDueDate := false

		var changes []string

		if cmd.Flags().Changed("content") {
			content = &todoEditContent
			changes = append(changes, "Updated content to \""+todoEditContent+"\"")
		}

		if cmd.Flags().Changed("due") {
			if todoEditDue == "" {
				clearDueDate = true
				changes = append(changes, "Removed due date")
			} else {
				parsed, err := dateparse.ParseDate(todoEditDue)
				if err != nil {
					return err
				}
				dueDate = &parsed
				changes = append(changes, "due date to "+parsed.Format("2006-01-02"))
			}
		}

		if len(todoEditTags) > 0 {
			changes = append(changes, "tags to "+strings.Join(formatTags(todoEditTags), " "))
		}

		if err := repository.UpdateTodo(database.DB, id, content, todoEditTags, dueDate, clearDueDate); err != nil {
			return err
		}

		if len(changes) > 0 {
			newTodo, err := repository.GetTodoByID(database.DB, id)
			if err != nil {
				return err
			}
			_ = oldTodo
			if err := activity.LogTodoUpdated(database.DB, newTodo, changes); err != nil {
				return err
			}
		}

		return nil
	},
}

func formatTags(tags []string) []string {
	formatted := make([]string, len(tags))
	for i, tag := range tags {
		formatted[i] = "#" + tag
	}
	return formatted
}

func init() {
	todoEditCmd.Flags().StringVar(&todoEditContent, "content", "", "New content for the todo")
	todoEditCmd.Flags().StringSliceVar(&todoEditTags, "tag", []string{}, "Replace tags")
	todoEditCmd.Flags().StringVar(&todoEditDue, "due", "", "Due date (YYYY-MM-DD or natural language)")
}

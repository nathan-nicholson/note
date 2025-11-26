package cmd

import (
	"time"

	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/dateparse"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var (
	todoAddTags []string
	todoAddDue  string
)

var todoAddCmd = &cobra.Command{
	Use:   "add <content>",
	Short: "Create a new todo",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content := args[0]

		activeProject, err := repository.GetActiveProject(database.DB)
		if err != nil {
			return err
		}

		tags := append(todoAddTags, activeProject.Name)

		var dueDate *time.Time
		if todoAddDue != "" {
			parsed, err := dateparse.ParseDate(todoAddDue)
			if err != nil {
				return err
			}
			dueDate = &parsed
		}

		todo, err := repository.CreateTodo(database.DB, content, tags, dueDate)
		if err != nil {
			return err
		}

		if err := repository.UpdateProjectLastActivity(database.DB, activeProject.ID); err != nil {
			return err
		}

		if err := activity.LogTodoCreated(database.DB, todo); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	todoAddCmd.Flags().StringSliceVar(&todoAddTags, "tag", []string{}, "Tags for the todo")
	todoAddCmd.Flags().StringVar(&todoAddDue, "due", "", "Due date (YYYY-MM-DD or natural language)")
}

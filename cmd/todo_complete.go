package cmd

import (
	"strconv"

	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var todoCompleteCmd = &cobra.Command{
	Use:   "complete <id>",
	Short: "Mark a todo as complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		todo, err := repository.GetTodoByID(database.DB, id)
		if err != nil {
			return err
		}

		if err := repository.CompleteTodo(database.DB, id); err != nil {
			return err
		}

		if err := activity.LogTodoCompleted(database.DB, todo); err != nil {
			return err
		}

		return nil
	},
}

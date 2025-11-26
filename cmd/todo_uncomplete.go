package cmd

import (
	"strconv"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var todoUncompleteCmd = &cobra.Command{
	Use:   "uncomplete <id>",
	Short: "Mark a todo as incomplete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		return repository.UncompleteTodo(database.DB, id)
	},
}

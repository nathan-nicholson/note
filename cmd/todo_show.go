package cmd

import (
	"fmt"
	"strconv"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/display"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var todoShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a todo",
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

		fmt.Println(display.FormatTodo(todo))
		return nil
	},
}

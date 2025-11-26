package cmd

import (
	"fmt"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/display"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var (
	todoListComplete   bool
	todoListIncomplete bool
	todoListTags       []string
	todoListOverdue    bool
)

var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List todos",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := repository.TodoListOptions{
			Complete:   todoListComplete,
			Incomplete: todoListIncomplete,
			Tags:       todoListTags,
			Overdue:    todoListOverdue,
		}

		todos, err := repository.ListTodos(database.DB, opts)
		if err != nil {
			return err
		}

		output := display.FormatTodoList(todos)
		if output != "" {
			fmt.Println(output)
		}

		return nil
	},
}

func init() {
	todoListCmd.Flags().BoolVar(&todoListComplete, "complete", false, "Show only complete todos")
	todoListCmd.Flags().BoolVar(&todoListIncomplete, "incomplete", false, "Show only incomplete todos")
	todoListCmd.Flags().StringSliceVar(&todoListTags, "tag", []string{}, "Filter by tags")
	todoListCmd.Flags().BoolVar(&todoListOverdue, "overdue", false, "Show only overdue todos")
}

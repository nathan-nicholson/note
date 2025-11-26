package cmd

import (
	"github.com/spf13/cobra"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage todos",
	Long:  `Create, list, and manage actionable todo items.`,
}

func init() {
	todoCmd.AddCommand(todoAddCmd)
	todoCmd.AddCommand(todoListCmd)
	todoCmd.AddCommand(todoEditCmd)
	todoCmd.AddCommand(todoDeleteCmd)
	todoCmd.AddCommand(todoShowCmd)
	todoCmd.AddCommand(todoCompleteCmd)
	todoCmd.AddCommand(todoUncompleteCmd)

	todoCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return todoAddCmd.RunE(cmd, args)
		}
		return cmd.Help()
	}
	todoCmd.Flags().StringSliceVar(&todoAddTags, "tag", []string{}, "Tags for the todo")
	todoCmd.Flags().StringVar(&todoAddDue, "due", "", "Due date (YYYY-MM-DD or natural language)")
}

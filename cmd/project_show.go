package cmd

import (
	"fmt"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/display"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectShowCmd = &cobra.Command{
	Use:   "show <project-name>",
	Short: "Show project details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		project, err := repository.GetProjectByName(database.DB, projectName)
		if err != nil {
			return err
		}

		fmt.Println(display.FormatProject(project))
		return nil
	},
}

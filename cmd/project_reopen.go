package cmd

import (
	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectReopenCmd = &cobra.Command{
	Use:   "reopen <project-name>",
	Short: "Reopen a closed project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		project, err := repository.GetProjectByName(database.DB, projectName)
		if err != nil {
			return err
		}

		if err := repository.ReopenProject(database.DB, project.ID); err != nil {
			return err
		}

		if err := activity.LogProjectReopened(database.DB, projectName); err != nil {
			return err
		}

		return nil
	},
}

package cmd

import (
	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectDeleteCmd = &cobra.Command{
	Use:   "delete <project-name>",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		project, err := repository.GetProjectByName(database.DB, projectName)
		if err != nil {
			return err
		}

		activeProject, err := repository.GetActiveProject(database.DB)
		if err != nil {
			return err
		}

		if activeProject.Name == projectName {
			if err := activity.LogProjectDeactivated(database.DB, projectName); err != nil {
				return err
			}

			homeProject, err := repository.GetProjectByName(database.DB, "home")
			if err != nil {
				return err
			}

			if err := repository.SetActiveProject(database.DB, homeProject.ID); err != nil {
				return err
			}

			if err := activity.LogProjectActivated(database.DB, "home"); err != nil {
				return err
			}
		}

		if err := activity.LogProjectDeleted(database.DB, project); err != nil {
			return err
		}

		if err := repository.DeleteProject(database.DB, projectName); err != nil {
			return err
		}

		return nil
	},
}

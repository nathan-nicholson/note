package cmd

import (
	"fmt"

	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project [project-name]",
	Short: "Manage projects",
	Long:  `Create, list, and switch between projects.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		projectName := args[0]

		project, err := repository.GetProjectByName(database.DB, projectName)
		if err != nil {
			return err
		}

		if project.IsClosed {
			return fmt.Errorf("Cannot set closed project '%s' as active. Reopen it with: note project reopen %s", projectName, projectName)
		}

		currentActive, err := repository.GetActiveProject(database.DB)
		if err != nil {
			return err
		}

		if currentActive.Name == projectName {
			return nil
		}

		if err := activity.LogProjectDeactivated(database.DB, currentActive.Name); err != nil {
			return err
		}

		if err := repository.SetActiveProject(database.DB, project.ID); err != nil {
			return err
		}

		if err := activity.LogProjectActivated(database.DB, projectName); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectCloseCmd)
	projectCmd.AddCommand(projectReopenCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectStatusCmd)
	projectCmd.AddCommand(projectShowCmd)
	projectCmd.AddCommand(projectEditCmd)
	projectCmd.AddCommand(projectDeleteCmd)
}

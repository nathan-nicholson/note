package cmd

import (
	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectCreateTags []string

var projectCreateCmd = &cobra.Command{
	Use:   "create <project-name>",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		project, err := repository.CreateProject(database.DB, projectName, projectCreateTags)
		if err != nil {
			return err
		}

		if err := activity.LogProjectCreated(database.DB, project); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	projectCreateCmd.Flags().StringSliceVar(&projectCreateTags, "tag", []string{}, "Tags for the project")
}

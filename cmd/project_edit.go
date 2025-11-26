package cmd

import (
	"strings"

	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectEditTags []string

var projectEditCmd = &cobra.Command{
	Use:   "edit <project-name>",
	Short: "Edit project tags",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		project, err := repository.GetProjectByName(database.DB, projectName)
		if err != nil {
			return err
		}

		if len(projectEditTags) > 0 {
			if err := repository.UpdateProjectTags(database.DB, project.ID, projectEditTags); err != nil {
				return err
			}

			formattedTags := make([]string, len(projectEditTags))
			for i, tag := range projectEditTags {
				formattedTags[i] = "#" + tag
			}
			changes := "Updated tags to " + strings.Join(formattedTags, " ")

			if err := activity.LogProjectUpdated(database.DB, projectName, changes); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	projectEditCmd.Flags().StringSliceVar(&projectEditTags, "tag", []string{}, "Replace project tags")
}

package cmd

import (
	"fmt"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/display"
	"github.com/nathan-nicholson/note/internal/models"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectStatusAll bool

var projectStatusCmd = &cobra.Command{
	Use:   "status [project-name]",
	Short: "Show project status",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var project *models.Project
		var err error

		if len(args) == 0 {
			project, err = repository.GetActiveProject(database.DB)
			if err != nil {
				return err
			}
		} else {
			project, err = repository.GetProjectByName(database.DB, args[0])
			if err != nil {
				return err
			}
		}

		output, err := display.FormatProjectStatus(database.DB, project, projectStatusAll)
		if err != nil {
			return err
		}

		fmt.Println(output)
		return nil
	},
}

func init() {
	projectStatusCmd.Flags().BoolVar(&projectStatusAll, "all", false, "Show all tasks including completed")
}

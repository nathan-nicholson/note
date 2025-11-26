package cmd

import (
	"fmt"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/display"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectListAll bool

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := repository.ListProjects(database.DB, projectListAll)
		if err != nil {
			return err
		}

		output, err := display.FormatProjectList(database.DB, projects, projectListAll)
		if err != nil {
			return err
		}

		if output != "" {
			fmt.Println(output)
		}

		return nil
	},
}

func init() {
	projectListCmd.Flags().BoolVar(&projectListAll, "all", false, "Include closed projects")
}

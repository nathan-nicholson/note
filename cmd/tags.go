package cmd

import (
	"fmt"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/display"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		tags, err := repository.ListAllTags(database.DB)
		if err != nil {
			return err
		}

		output, err := display.FormatTagList(database.DB, tags)
		if err != nil {
			return err
		}

		if output != "" {
			fmt.Println(output)
		}

		return nil
	},
}

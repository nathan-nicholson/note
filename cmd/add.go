package cmd

import (
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var (
	addTags      []string
	addImportant bool
)

var addCmd = &cobra.Command{
	Use:   "add <content>",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content := args[0]

		activeProject, err := repository.GetActiveProject(database.DB)
		if err != nil {
			return err
		}

		tags := append(addTags, activeProject.Name)

		_, err = repository.CreateNote(database.DB, content, tags, addImportant)
		if err != nil {
			return err
		}

		if err := repository.UpdateProjectLastActivity(database.DB, activeProject.ID); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringSliceVar(&addTags, "tag", []string{}, "Tags for the note")
	addCmd.Flags().BoolVar(&addImportant, "important", false, "Mark note as important")
}

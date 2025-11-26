package cmd

import (
	"strconv"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var (
	editContent   string
	editTags      []string
	editImportant bool
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		var content *string
		var important *bool

		if cmd.Flags().Changed("content") {
			content = &editContent
		}

		if cmd.Flags().Changed("important") {
			important = &editImportant
		}

		return repository.UpdateNote(database.DB, id, content, editTags, important)
	},
}

func init() {
	editCmd.Flags().StringVar(&editContent, "content", "", "New content for the note")
	editCmd.Flags().StringSliceVar(&editTags, "tag", []string{}, "Replace tags")
	editCmd.Flags().BoolVar(&editImportant, "important", false, "Mark as important")
}

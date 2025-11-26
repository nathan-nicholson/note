package cmd

import (
	"fmt"
	"strconv"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/display"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		note, err := repository.GetNoteByID(database.DB, id)
		if err != nil {
			return err
		}

		fmt.Println(display.FormatNote(note))
		return nil
	},
}

package cmd

import (
	"fmt"
	"time"

	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/dateparse"
	"github.com/nathan-nicholson/note/internal/display"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var (
	listStart     string
	listEnd       string
	listTags      []string
	listImportant bool
	listShowIDs   bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := repository.NoteListOptions{
			Tags:      listTags,
			Important: listImportant,
		}

		if listStart == "" && listEnd == "" {
			today := time.Now()
			opts.StartDate = &today
			opts.EndDate = &today
		} else {
			if listStart != "" {
				start, err := dateparse.ParseDate(listStart)
				if err != nil {
					return err
				}
				opts.StartDate = &start
			}

			if listEnd != "" {
				end, err := dateparse.ParseDate(listEnd)
				if err != nil {
					return err
				}
				opts.EndDate = &end
			}
		}

		notes, err := repository.ListNotes(database.DB, opts)
		if err != nil {
			return err
		}

		output := display.FormatNoteList(notes, listShowIDs)
		if output != "" {
			fmt.Println(output)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&listStart, "start", "", "Start date (YYYY-MM-DD or natural language)")
	listCmd.Flags().StringVar(&listEnd, "end", "", "End date (YYYY-MM-DD or natural language)")
	listCmd.Flags().StringSliceVar(&listTags, "tag", []string{}, "Filter by tags")
	listCmd.Flags().BoolVar(&listImportant, "important", false, "Show only important notes")
	listCmd.Flags().BoolVar(&listShowIDs, "show-ids", false, "Show note IDs")
}

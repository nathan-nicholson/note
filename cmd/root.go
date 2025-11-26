package cmd

import (
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var (
	rootTags      []string
	rootImportant bool
)

var rootCmd = &cobra.Command{
	Use:   "note [content]",
	Short: "A lightweight CLI tool for capturing notes and managing todos",
	Long:  `note is a fast, keyboard-driven tool for capturing thoughts and tasks with project-based organization.`,
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		content := args[0]

		activeProject, err := repository.GetActiveProject(database.DB)
		if err != nil {
			return err
		}

		tags := append(rootTags, activeProject.Name)

		_, err = repository.CreateNote(database.DB, content, tags, rootImportant)
		if err != nil {
			return err
		}

		if err := repository.UpdateProjectLastActivity(database.DB, activeProject.ID); err != nil {
			return err
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringSliceVar(&rootTags, "tag", []string{}, "Tags for the note")
	rootCmd.Flags().BoolVar(&rootImportant, "important", false, "Mark note as important")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(todoCmd)
	rootCmd.AddCommand(projectCmd)
}

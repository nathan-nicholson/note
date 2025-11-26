package cmd

import (
	"fmt"
	"runtime"

	"github.com/nathan-nicholson/note/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display the current version of note along with build information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("note version %s\n", version.Version)
		fmt.Printf("Built with %s for %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

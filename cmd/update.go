package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nathan-nicholson/note/internal/update"
	"github.com/spf13/cobra"
)

var (
	updateCheck bool
	updateForce bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Long:  `Check for the latest version of note and optionally upgrade.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check for updates
		fmt.Println("Checking for updates...")
		updateInfo, err := update.CheckForUpdates()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		// Display current and latest versions
		fmt.Printf("\nCurrent version: %s\n", color.CyanString(updateInfo.CurrentVersion))
		fmt.Printf("Latest version:  %s\n", color.GreenString(updateInfo.LatestVersion))

		// Check if update is available
		if !updateInfo.IsNewer {
			fmt.Println(color.GreenString("\n✓ You're already on the latest version!"))
			return nil
		}

		// Show that update is available
		fmt.Println(color.YellowString("\n⚠ A new version is available!"))

		// Show release notes if available
		if updateInfo.ReleaseNotes != "" {
			fmt.Println("\nRelease notes:")
			fmt.Println(strings.Repeat("-", 60))
			// Show first few lines of release notes
			lines := strings.Split(updateInfo.ReleaseNotes, "\n")
			maxLines := 10
			for i, line := range lines {
				if i >= maxLines {
					fmt.Println("...")
					fmt.Printf("View full release notes: %s\n", updateInfo.UpdateURL)
					break
				}
				fmt.Println(line)
			}
			fmt.Println(strings.Repeat("-", 60))
		}

		// If check-only mode, show instructions and exit
		if updateCheck {
			fmt.Println()
			method := update.DetectInstallMethod()
			fmt.Println(update.GetUpdateInstructions(method))
			return nil
		}

		// Detect installation method
		method := update.DetectInstallMethod()
		canAutoUpdate := method == update.InstallHomebrew || method == update.InstallGoInstall

		if !canAutoUpdate {
			fmt.Println("\nAutomatic update is not available for your installation method.")
			fmt.Println(update.GetUpdateInstructions(method))
			return nil
		}

		// Prompt user for confirmation
		if !updateForce {
			fmt.Printf("\nWould you like to update now? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("\nUpdate cancelled. Run 'note update' again when you're ready.")
				return nil
			}
		}

		// Perform the update
		fmt.Println("\nUpdating note...")
		if err := update.PerformUpdate(method); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}

		fmt.Println(color.GreenString("\n✓ Successfully updated to version %s!", updateInfo.LatestVersion))
		fmt.Println("\nPlease restart note to use the new version.")

		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateCheck, "check", false, "Only check for updates, don't install")
	updateCmd.Flags().BoolVarP(&updateForce, "yes", "y", false, "Skip confirmation prompt")
	rootCmd.AddCommand(updateCmd)
}

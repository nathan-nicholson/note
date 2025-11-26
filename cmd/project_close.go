package cmd

import (
	"fmt"

	"github.com/nathan-nicholson/note/internal/activity"
	"github.com/nathan-nicholson/note/internal/database"
	"github.com/nathan-nicholson/note/internal/repository"
	"github.com/spf13/cobra"
)

var projectCloseCmd = &cobra.Command{
	Use:   "close <project-name>",
	Short: "Close a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		project, err := repository.GetProjectByName(database.DB, projectName)
		if err != nil {
			return err
		}

		incompleteTodos, err := repository.GetIncompleteTodosForProject(database.DB, projectName)
		if err != nil {
			return err
		}

		if len(incompleteTodos) > 0 {
			fmt.Printf("Error: Cannot close project '%s' - %d incomplete todos remaining\n\n", projectName, len(incompleteTodos))
			fmt.Println("Incomplete Tasks:")
			for _, todo := range incompleteTodos {
				fmt.Printf("  [ ] [#%d] ", todo.ID)
				if todo.DueDate.Valid {
					fmt.Printf("%s  ", todo.DueDate.Time.Format("2006-01-02"))
				}
				fmt.Print(todo.Content)
				if len(todo.Tags) > 0 {
					fmt.Print(" ")
					for _, tag := range todo.Tags {
						fmt.Printf("#%s ", tag)
					}
				}
				fmt.Println()
			}
			fmt.Println("\nComplete all todos before closing the project.")
			return fmt.Errorf("cannot close project with incomplete todos")
		}

		openCount, err := repository.CountOpenProjects(database.DB)
		if err != nil {
			return err
		}

		if openCount == 1 && projectName == "home" {
			return fmt.Errorf("Cannot close 'home' project - it is the only open project. Create or reopen another project first.")
		}

		activeProject, err := repository.GetActiveProject(database.DB)
		if err != nil {
			return err
		}

		if activeProject.Name == projectName {
			if err := activity.LogProjectDeactivated(database.DB, projectName); err != nil {
				return err
			}

			homeProject, err := repository.GetProjectByName(database.DB, "home")
			if err != nil {
				openProjects, err := repository.ListProjects(database.DB, false)
				if err != nil {
					return err
				}

				for _, p := range openProjects {
					if p.Name != projectName {
						if err := repository.SetActiveProject(database.DB, p.ID); err != nil {
							return err
						}
						if err := activity.LogProjectActivated(database.DB, p.Name); err != nil {
							return err
						}
						break
					}
				}
			} else if !homeProject.IsClosed {
				if err := repository.SetActiveProject(database.DB, homeProject.ID); err != nil {
					return err
				}
				if err := activity.LogProjectActivated(database.DB, "home"); err != nil {
					return err
				}
			}
		}

		if err := repository.CloseProject(database.DB, project.ID); err != nil {
			return err
		}

		if err := activity.LogProjectClosed(database.DB, projectName); err != nil {
			return err
		}

		fmt.Printf("Project '%s' closed successfully.\n", projectName)
		return nil
	},
}

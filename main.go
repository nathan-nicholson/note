package main

import (
	"fmt"
	"os"

	"github.com/nathan-nicholson/note/cmd"
	"github.com/nathan-nicholson/note/internal/database"
)

func main() {
	if err := database.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer database.CloseDB()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

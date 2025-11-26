package display

import (
	"fmt"
	"strings"

	"github.com/nathan-nicholson/note/internal/models"
)

func FormatNoteList(notes []models.Note, showIDs bool) string {
	if len(notes) == 0 {
		return ""
	}

	var output strings.Builder
	currentDate := ""

	for _, note := range notes {
		noteDate := note.CreatedAt.Format("2006-01-02")

		if noteDate != currentDate {
			if currentDate != "" {
				output.WriteString("\n")
			}
			output.WriteString(noteDate + "\n")
			currentDate = noteDate
		}

		output.WriteString("  ")

		if showIDs {
			output.WriteString(fmt.Sprintf("[#%d] ", note.ID))
		}

		output.WriteString(note.CreatedAt.Format("03:04 PM"))
		output.WriteString("  ")

		if note.IsImportant {
			output.WriteString("[!] ")
		}

		output.WriteString(note.Content)

		if len(note.Tags) > 0 {
			output.WriteString(" ")
			for _, tag := range note.Tags {
				output.WriteString("#" + tag + " ")
			}
		}

		output.WriteString("\n")
	}

	return strings.TrimSpace(output.String())
}

func FormatNote(note *models.Note) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("ID: %d\n", note.ID))
	output.WriteString(fmt.Sprintf("Created: %s\n", note.CreatedAt.Format("2006-01-02 03:04 PM")))
	output.WriteString(fmt.Sprintf("Updated: %s\n", note.UpdatedAt.Format("2006-01-02 03:04 PM")))

	if note.IsImportant {
		output.WriteString("Important: Yes\n")
	}

	if len(note.Tags) > 0 {
		output.WriteString("Tags: ")
		for i, tag := range note.Tags {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString("#" + tag)
		}
		output.WriteString("\n")
	}

	output.WriteString("\n")
	output.WriteString(note.Content)

	return output.String()
}


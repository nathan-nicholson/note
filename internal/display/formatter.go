package display

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/nathan-nicholson/note/internal/models"
)

func FormatTagList(db *sql.DB, tags []models.Tag) (string, error) {
	if len(tags) == 0 {
		return "", nil
	}

	var output strings.Builder

	for _, tag := range tags {
		var usageCount int
		err := db.QueryRow(`
			SELECT
				(SELECT COUNT(*) FROM note_tags WHERE tag_id = ?) +
				(SELECT COUNT(*) FROM todo_tags WHERE tag_id = ?)
		`, tag.ID, tag.ID).Scan(&usageCount)

		if err != nil {
			return "", err
		}

		output.WriteString(fmt.Sprintf("%s (%d uses)\n", tag.Name, usageCount))
	}

	return strings.TrimSpace(output.String()), nil
}

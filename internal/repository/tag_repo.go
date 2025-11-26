package repository

import (
	"database/sql"
	"github.com/nathan-nicholson/note/internal/models"
)

func GetOrCreateTag(db *sql.DB, name string) (int, error) {
	var tagID int
	err := db.QueryRow("SELECT id FROM tags WHERE name = ?", name).Scan(&tagID)
	if err == nil {
		return tagID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	result, err := db.Exec("INSERT INTO tags (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetTagsForNote(db *sql.DB, noteID int) ([]string, error) {
	rows, err := db.Query(`
		SELECT t.name
		FROM tags t
		JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = ?
		ORDER BY t.name
	`, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func GetTagsForTodo(db *sql.DB, todoID int) ([]string, error) {
	rows, err := db.Query(`
		SELECT t.name
		FROM tags t
		JOIN todo_tags tt ON t.id = tt.tag_id
		WHERE tt.todo_id = ?
		ORDER BY t.name
	`, todoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func GetTagsForProject(db *sql.DB, projectID int) ([]string, error) {
	rows, err := db.Query(`
		SELECT t.name
		FROM tags t
		JOIN project_tags pt ON t.id = pt.tag_id
		WHERE pt.project_id = ?
		ORDER BY t.name
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func AddTagsToNote(db *sql.DB, noteID int, tags []string) error {
	for _, tagName := range tags {
		tagID, err := GetOrCreateTag(db, tagName)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT OR IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)", noteID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddTagsToTodo(db *sql.DB, todoID int, tags []string) error {
	for _, tagName := range tags {
		tagID, err := GetOrCreateTag(db, tagName)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT OR IGNORE INTO todo_tags (todo_id, tag_id) VALUES (?, ?)", todoID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddTagsToProject(db *sql.DB, projectID int, tags []string) error {
	for _, tagName := range tags {
		tagID, err := GetOrCreateTag(db, tagName)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT OR IGNORE INTO project_tags (project_id, tag_id) VALUES (?, ?)", projectID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReplaceNoteTags(db *sql.DB, noteID int, tags []string) error {
	_, err := db.Exec("DELETE FROM note_tags WHERE note_id = ?", noteID)
	if err != nil {
		return err
	}
	return AddTagsToNote(db, noteID, tags)
}

func ReplaceTodoTags(db *sql.DB, todoID int, tags []string) error {
	_, err := db.Exec("DELETE FROM todo_tags WHERE todo_id = ?", todoID)
	if err != nil {
		return err
	}
	return AddTagsToTodo(db, todoID, tags)
}

func ReplaceProjectTags(db *sql.DB, projectID int, tags []string) error {
	_, err := db.Exec("DELETE FROM project_tags WHERE project_id = ?", projectID)
	if err != nil {
		return err
	}
	return AddTagsToProject(db, projectID, tags)
}

func ListAllTags(db *sql.DB) ([]models.Tag, error) {
	rows, err := db.Query(`
		SELECT t.id, t.name,
			(SELECT COUNT(*) FROM note_tags WHERE tag_id = t.id) +
			(SELECT COUNT(*) FROM todo_tags WHERE tag_id = t.id) AS usage_count
		FROM tags t
		ORDER BY usage_count DESC, t.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		var usageCount int
		if err := rows.Scan(&tag.ID, &tag.Name, &usageCount); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

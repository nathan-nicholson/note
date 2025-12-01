package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/nathan-nicholson/note/internal/models"
)

func CreateNote(db *sql.DB, content string, tags []string, isImportant bool) (*models.Note, error) {
	now := time.Now()
	result, err := db.Exec(`
		INSERT INTO notes (content, is_important, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, content, isImportant, now, now)
	if err != nil {
		return nil, err
	}

	noteID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := AddTagsToNote(db, int(noteID), tags); err != nil {
		return nil, err
	}

	return GetNoteByID(db, int(noteID))
}

func GetNoteByID(db *sql.DB, id int) (*models.Note, error) {
	var note models.Note
	err := db.QueryRow(`
		SELECT id, content, created_at, updated_at, is_important
		FROM notes
		WHERE id = ?
	`, id).Scan(&note.ID, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.IsImportant)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Note #%d not found", id)
		}
		return nil, err
	}

	tags, err := GetTagsForNote(db, id)
	if err != nil {
		return nil, err
	}
	note.Tags = tags

	return &note, nil
}

type NoteListOptions struct {
	StartDate   *time.Time
	EndDate     *time.Time
	Tags        []string
	Important   bool
}

func ListNotes(db *sql.DB, opts NoteListOptions) ([]models.Note, error) {
	query := `
		SELECT DISTINCT n.id, n.content, n.created_at, n.updated_at, n.is_important
		FROM notes n
	`

	var conditions []string
	var args []interface{}

	if len(opts.Tags) > 0 {
		query += `
			JOIN note_tags nt ON n.id = nt.note_id
			JOIN tags t ON nt.tag_id = t.id
		`
		placeholders := make([]string, len(opts.Tags))
		for i, tag := range opts.Tags {
			placeholders[i] = "?"
			args = append(args, tag)
		}
		conditions = append(conditions, fmt.Sprintf("t.name IN (%s)", strings.Join(placeholders, ",")))
	}

	if opts.StartDate != nil {
		conditions = append(conditions, "DATE(n.created_at) >= DATE(?)")
		args = append(args, opts.StartDate.Format("2006-01-02"))
	}

	if opts.EndDate != nil {
		conditions = append(conditions, "DATE(n.created_at) <= DATE(?)")
		args = append(args, opts.EndDate.Format("2006-01-02"))
	}

	if opts.Important {
		conditions = append(conditions, "n.is_important = 1")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if len(opts.Tags) > 0 {
		query += fmt.Sprintf(" GROUP BY n.id HAVING COUNT(DISTINCT t.name) = %d", len(opts.Tags))
	}

	query += " ORDER BY n.created_at ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.IsImportant); err != nil {
			return nil, err
		}

		tags, err := GetTagsForNote(db, note.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = tags

		notes = append(notes, note)
	}

	return notes, rows.Err()
}

func UpdateNote(db *sql.DB, id int, content *string, tags []string, isImportant *bool) error {
	now := time.Now()

	if content != nil {
		_, err := db.Exec(`
			UPDATE notes
			SET content = ?, updated_at = ?
			WHERE id = ?
		`, *content, now, id)
		if err != nil {
			return err
		}
	}

	if isImportant != nil {
		_, err := db.Exec(`
			UPDATE notes
			SET is_important = ?, updated_at = ?
			WHERE id = ?
		`, *isImportant, now, id)
		if err != nil {
			return err
		}
	}

	if len(tags) > 0 {
		if err := ReplaceNoteTags(db, id, tags); err != nil {
			return err
		}
		_, err := db.Exec("UPDATE notes SET updated_at = ? WHERE id = ?", now, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteNote(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Note #%d not found", id)
	}

	return nil
}

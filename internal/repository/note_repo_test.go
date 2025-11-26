package repository

import (
	"sort"
	"testing"
	"time"
)

func TestCreateNote(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name        string
		content     string
		tags        []string
		isImportant bool
	}{
		{
			name:        "simple note without tags",
			content:     "This is a test note",
			tags:        []string{},
			isImportant: false,
		},
		{
			name:        "important note with tags",
			content:     "This is important",
			tags:        []string{"work", "urgent"},
			isImportant: true,
		},
		{
			name:        "note with multiple tags",
			content:     "Meeting notes",
			tags:        []string{"meeting", "team", "project"},
			isImportant: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note, err := CreateNote(db, tt.content, tt.tags, tt.isImportant)
			if err != nil {
				t.Fatalf("CreateNote() error = %v", err)
			}

			if note.ID == 0 {
				t.Error("CreateNote() returned note with ID = 0")
			}

			if note.Content != tt.content {
				t.Errorf("CreateNote() content = %q, want %q", note.Content, tt.content)
			}

			if note.IsImportant != tt.isImportant {
				t.Errorf("CreateNote() isImportant = %v, want %v", note.IsImportant, tt.isImportant)
			}

			if len(note.Tags) != len(tt.tags) {
				t.Errorf("CreateNote() tags count = %d, want %d", len(note.Tags), len(tt.tags))
			}

			expectedTags := make([]string, len(tt.tags))
			copy(expectedTags, tt.tags)
			sort.Strings(expectedTags)

			for i, tag := range note.Tags {
				if tag != expectedTags[i] {
					t.Errorf("CreateNote() tag[%d] = %q, want %q", i, tag, expectedTags[i])
				}
			}

			if note.CreatedAt.IsZero() {
				t.Error("CreateNote() CreatedAt is zero")
			}

			if note.UpdatedAt.IsZero() {
				t.Error("CreateNote() UpdatedAt is zero")
			}
		})
	}
}

func TestGetNoteByID(t *testing.T) {
	db := setupTestDB(t)

	createdNote, err := CreateNote(db, "Test note", []string{"test"}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	t.Run("existing note", func(t *testing.T) {
		note, err := GetNoteByID(db, createdNote.ID)
		if err != nil {
			t.Fatalf("GetNoteByID() error = %v", err)
		}

		if note.ID != createdNote.ID {
			t.Errorf("GetNoteByID() ID = %d, want %d", note.ID, createdNote.ID)
		}

		if note.Content != createdNote.Content {
			t.Errorf("GetNoteByID() Content = %q, want %q", note.Content, createdNote.Content)
		}
	})

	t.Run("non-existing note", func(t *testing.T) {
		_, err := GetNoteByID(db, 99999)
		if err == nil {
			t.Error("GetNoteByID() expected error for non-existing note, got nil")
		}
	})
}

func TestListNotes(t *testing.T) {
	db := setupTestDB(t)

	yesterday := time.Now().AddDate(0, 0, -1)
	today := time.Now()

	_, err := CreateNote(db, "Old note", []string{"old"}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	_, err = CreateNote(db, "Important work note", []string{"work"}, true)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	_, err = CreateNote(db, "Regular work note", []string{"work", "project"}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	t.Run("list all notes", func(t *testing.T) {
		notes, err := ListNotes(db, NoteListOptions{})
		if err != nil {
			t.Fatalf("ListNotes() error = %v", err)
		}

		if len(notes) != 3 {
			t.Errorf("ListNotes() returned %d notes, want 3", len(notes))
		}
	})

	t.Run("filter by important", func(t *testing.T) {
		notes, err := ListNotes(db, NoteListOptions{Important: true})
		if err != nil {
			t.Fatalf("ListNotes() error = %v", err)
		}

		if len(notes) != 1 {
			t.Errorf("ListNotes() returned %d notes, want 1", len(notes))
		}

		if len(notes) > 0 && !notes[0].IsImportant {
			t.Error("ListNotes() returned non-important note when filtering by important")
		}
	})

	t.Run("filter by single tag", func(t *testing.T) {
		notes, err := ListNotes(db, NoteListOptions{Tags: []string{"work"}})
		if err != nil {
			t.Fatalf("ListNotes() error = %v", err)
		}

		if len(notes) != 2 {
			t.Errorf("ListNotes() returned %d notes, want 2", len(notes))
		}
	})

	t.Run("filter by multiple tags", func(t *testing.T) {
		notes, err := ListNotes(db, NoteListOptions{Tags: []string{"work", "project"}})
		if err != nil {
			t.Fatalf("ListNotes() error = %v", err)
		}

		if len(notes) != 1 {
			t.Errorf("ListNotes() returned %d notes, want 1", len(notes))
		}
	})

	t.Run("filter by start date", func(t *testing.T) {
		notes, err := ListNotes(db, NoteListOptions{StartDate: &today})
		if err != nil {
			t.Fatalf("ListNotes() error = %v", err)
		}

		if len(notes) != 3 {
			t.Errorf("ListNotes() returned %d notes, want 3", len(notes))
		}
	})

	t.Run("filter by end date", func(t *testing.T) {
		notes, err := ListNotes(db, NoteListOptions{EndDate: &yesterday})
		if err != nil {
			t.Fatalf("ListNotes() error = %v", err)
		}

		if len(notes) != 0 {
			t.Errorf("ListNotes() returned %d notes, want 0", len(notes))
		}
	})
}

func TestUpdateNote(t *testing.T) {
	db := setupTestDB(t)

	note, err := CreateNote(db, "Original content", []string{"original"}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	t.Run("update content", func(t *testing.T) {
		newContent := "Updated content"
		err := UpdateNote(db, note.ID, &newContent, nil, nil)
		if err != nil {
			t.Fatalf("UpdateNote() error = %v", err)
		}

		updated, err := GetNoteByID(db, note.ID)
		if err != nil {
			t.Fatalf("GetNoteByID() error = %v", err)
		}

		if updated.Content != newContent {
			t.Errorf("UpdateNote() content = %q, want %q", updated.Content, newContent)
		}
	})

	t.Run("update importance", func(t *testing.T) {
		important := true
		err := UpdateNote(db, note.ID, nil, nil, &important)
		if err != nil {
			t.Fatalf("UpdateNote() error = %v", err)
		}

		updated, err := GetNoteByID(db, note.ID)
		if err != nil {
			t.Fatalf("GetNoteByID() error = %v", err)
		}

		if !updated.IsImportant {
			t.Error("UpdateNote() failed to update importance")
		}
	})

	t.Run("update tags", func(t *testing.T) {
		newTags := []string{"new", "tags"}
		err := UpdateNote(db, note.ID, nil, newTags, nil)
		if err != nil {
			t.Fatalf("UpdateNote() error = %v", err)
		}

		updated, err := GetNoteByID(db, note.ID)
		if err != nil {
			t.Fatalf("GetNoteByID() error = %v", err)
		}

		if len(updated.Tags) != len(newTags) {
			t.Errorf("UpdateNote() tags count = %d, want %d", len(updated.Tags), len(newTags))
		}
	})
}

func TestDeleteNote(t *testing.T) {
	db := setupTestDB(t)

	note, err := CreateNote(db, "Note to delete", []string{"delete"}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	t.Run("delete existing note", func(t *testing.T) {
		err := DeleteNote(db, note.ID)
		if err != nil {
			t.Fatalf("DeleteNote() error = %v", err)
		}

		_, err = GetNoteByID(db, note.ID)
		if err == nil {
			t.Error("DeleteNote() did not delete the note")
		}
	})

	t.Run("delete non-existing note", func(t *testing.T) {
		err := DeleteNote(db, 99999)
		if err == nil {
			t.Error("DeleteNote() expected error for non-existing note, got nil")
		}
	})
}

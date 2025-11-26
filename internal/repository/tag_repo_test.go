package repository

import (
	"testing"
)

func TestGetOrCreateTag(t *testing.T) {
	db := setupTestDB(t)

	t.Run("create new tag", func(t *testing.T) {
		tagID, err := GetOrCreateTag(db, "newtag")
		if err != nil {
			t.Fatalf("GetOrCreateTag() error = %v", err)
		}

		if tagID == 0 {
			t.Error("GetOrCreateTag() returned ID = 0")
		}
	})

	t.Run("get existing tag", func(t *testing.T) {
		firstID, err := GetOrCreateTag(db, "existingtag")
		if err != nil {
			t.Fatalf("GetOrCreateTag() first call error = %v", err)
		}

		secondID, err := GetOrCreateTag(db, "existingtag")
		if err != nil {
			t.Fatalf("GetOrCreateTag() second call error = %v", err)
		}

		if firstID != secondID {
			t.Errorf("GetOrCreateTag() returned different IDs: %d and %d", firstID, secondID)
		}
	})
}

func TestAddTagsToNote(t *testing.T) {
	db := setupTestDB(t)

	note, err := CreateNote(db, "Test note", []string{}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tags := []string{"tag1", "tag2", "tag3"}
	err = AddTagsToNote(db, note.ID, tags)
	if err != nil {
		t.Fatalf("AddTagsToNote() error = %v", err)
	}

	retrievedTags, err := GetTagsForNote(db, note.ID)
	if err != nil {
		t.Fatalf("GetTagsForNote() error = %v", err)
	}

	if len(retrievedTags) != len(tags) {
		t.Errorf("GetTagsForNote() returned %d tags, want %d", len(retrievedTags), len(tags))
	}
}

func TestReplaceNoteTags(t *testing.T) {
	db := setupTestDB(t)

	note, err := CreateNote(db, "Test note", []string{"old1", "old2"}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	newTags := []string{"new1", "new2", "new3"}
	err = ReplaceNoteTags(db, note.ID, newTags)
	if err != nil {
		t.Fatalf("ReplaceNoteTags() error = %v", err)
	}

	retrievedTags, err := GetTagsForNote(db, note.ID)
	if err != nil {
		t.Fatalf("GetTagsForNote() error = %v", err)
	}

	if len(retrievedTags) != len(newTags) {
		t.Errorf("GetTagsForNote() returned %d tags, want %d", len(retrievedTags), len(newTags))
	}

	for _, newTag := range newTags {
		found := false
		for _, retrievedTag := range retrievedTags {
			if retrievedTag == newTag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ReplaceNoteTags() tag %q not found in retrieved tags", newTag)
		}
	}
}

func TestListAllTags(t *testing.T) {
	db := setupTestDB(t)

	_, err := CreateNote(db, "Note 1", []string{"tag1", "tag2"}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	_, err = CreateNote(db, "Note 2", []string{"tag2", "tag3"}, false)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tags, err := ListAllTags(db)
	if err != nil {
		t.Fatalf("ListAllTags() error = %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("ListAllTags() returned %d tags, want 3", len(tags))
	}

	tagNames := make(map[string]bool)
	for _, tag := range tags {
		tagNames[tag.Name] = true
		if tag.ID == 0 {
			t.Errorf("ListAllTags() returned tag with ID = 0")
		}
	}

	expectedTags := []string{"tag1", "tag2", "tag3"}
	for _, expected := range expectedTags {
		if !tagNames[expected] {
			t.Errorf("ListAllTags() missing expected tag %q", expected)
		}
	}
}

func TestAddTagsToTodo(t *testing.T) {
	db := setupTestDB(t)

	todo, err := CreateTodo(db, "Test todo", []string{}, nil)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tags := []string{"urgent", "work"}
	err = AddTagsToTodo(db, todo.ID, tags)
	if err != nil {
		t.Fatalf("AddTagsToTodo() error = %v", err)
	}

	retrievedTags, err := GetTagsForTodo(db, todo.ID)
	if err != nil {
		t.Fatalf("GetTagsForTodo() error = %v", err)
	}

	if len(retrievedTags) != len(tags) {
		t.Errorf("GetTagsForTodo() returned %d tags, want %d", len(retrievedTags), len(tags))
	}
}

func TestReplaceTodoTags(t *testing.T) {
	db := setupTestDB(t)

	todo, err := CreateTodo(db, "Test todo", []string{"old"}, nil)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	newTags := []string{"new1", "new2"}
	err = ReplaceTodoTags(db, todo.ID, newTags)
	if err != nil {
		t.Fatalf("ReplaceTodoTags() error = %v", err)
	}

	retrievedTags, err := GetTagsForTodo(db, todo.ID)
	if err != nil {
		t.Fatalf("GetTagsForTodo() error = %v", err)
	}

	if len(retrievedTags) != len(newTags) {
		t.Errorf("GetTagsForTodo() returned %d tags, want %d", len(retrievedTags), len(newTags))
	}
}

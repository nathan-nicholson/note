package repository

import (
	"testing"
	"time"
)

func TestNoteCreatedAtNotInFuture(t *testing.T) {
	db := setupTestDB(t)

	// Capture time before creating note
	beforeCreate := time.Now()

	// Create a note
	note, err := CreateNote(db, "Test note", []string{}, false)
	if err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	// Capture time after creating note
	afterCreate := time.Now()

	// The created_at should be between beforeCreate and afterCreate
	// It should NOT be in the future
	if note.CreatedAt.After(afterCreate) {
		t.Errorf("Note created_at is in the future!\n"+
			"  Before create: %v\n"+
			"  Created at:    %v\n"+
			"  After create:  %v\n"+
			"  Difference:    %v",
			beforeCreate, note.CreatedAt, afterCreate,
			note.CreatedAt.Sub(afterCreate))
	}

	if note.CreatedAt.Before(beforeCreate) {
		t.Errorf("Note created_at is before we called CreateNote!\n"+
			"  Before create: %v\n"+
			"  Created at:    %v\n"+
			"  Difference:    %v",
			beforeCreate, note.CreatedAt,
			beforeCreate.Sub(note.CreatedAt))
	}
}

func TestTodoCreatedAtNotInFuture(t *testing.T) {
	db := setupTestDB(t)

	// Capture time before creating todo
	beforeCreate := time.Now()

	// Create a todo
	todo, err := CreateTodo(db, "Test todo", []string{}, nil)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Capture time after creating todo
	afterCreate := time.Now()

	// The created_at should be between beforeCreate and afterCreate
	// It should NOT be in the future
	if todo.CreatedAt.After(afterCreate) {
		t.Errorf("Todo created_at is in the future!\n"+
			"  Before create: %v\n"+
			"  Created at:    %v\n"+
			"  After create:  %v\n"+
			"  Difference:    %v",
			beforeCreate, todo.CreatedAt, afterCreate,
			todo.CreatedAt.Sub(afterCreate))
	}

	if todo.CreatedAt.Before(beforeCreate) {
		t.Errorf("Todo created_at is before we called CreateTodo!\n"+
			"  Before create: %v\n"+
			"  Created at:    %v\n"+
			"  Difference:    %v",
			beforeCreate, todo.CreatedAt,
			beforeCreate.Sub(todo.CreatedAt))
	}
}

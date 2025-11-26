package repository

import (
	"testing"
	"time"
)

func TestCreateTodo(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name    string
		content string
		tags    []string
		dueDate *time.Time
	}{
		{
			name:    "simple todo without due date",
			content: "Buy groceries",
			tags:    []string{},
			dueDate: nil,
		},
		{
			name:    "todo with due date and tags",
			content: "Complete project",
			tags:    []string{"work", "urgent"},
			dueDate: timePtr(time.Now().AddDate(0, 0, 7)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := CreateTodo(db, tt.content, tt.tags, tt.dueDate)
			if err != nil {
				t.Fatalf("CreateTodo() error = %v", err)
			}

			if todo.ID == 0 {
				t.Error("CreateTodo() returned todo with ID = 0")
			}

			if todo.Content != tt.content {
				t.Errorf("CreateTodo() content = %q, want %q", todo.Content, tt.content)
			}

			if todo.IsComplete {
				t.Error("CreateTodo() new todo should not be complete")
			}

			if tt.dueDate != nil {
				if !todo.DueDate.Valid {
					t.Error("CreateTodo() due date is not valid, expected a date")
				} else {
					// Compare only the date parts (year, month, day)
					gotDate := todo.DueDate.Time.Format("2006-01-02")
					wantDate := tt.dueDate.Format("2006-01-02")
					if gotDate != wantDate {
						t.Errorf("CreateTodo() due date = %v, want %v", gotDate, wantDate)
					}
				}
			}

			if len(todo.Tags) != len(tt.tags) {
				t.Errorf("CreateTodo() tags count = %d, want %d", len(todo.Tags), len(tt.tags))
			}
		})
	}
}

func TestCompleteTodo(t *testing.T) {
	db := setupTestDB(t)

	todo, err := CreateTodo(db, "Test todo", []string{}, nil)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	err = CompleteTodo(db, todo.ID)
	if err != nil {
		t.Fatalf("CompleteTodo() error = %v", err)
	}

	updated, err := GetTodoByID(db, todo.ID)
	if err != nil {
		t.Fatalf("GetTodoByID() error = %v", err)
	}

	if !updated.IsComplete {
		t.Error("CompleteTodo() did not mark todo as complete")
	}

	if !updated.CompletedAt.Valid {
		t.Error("CompleteTodo() did not set completed_at timestamp")
	}
}

func TestUncompleteTodo(t *testing.T) {
	db := setupTestDB(t)

	todo, err := CreateTodo(db, "Test todo", []string{}, nil)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	if err := CompleteTodo(db, todo.ID); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	err = UncompleteTodo(db, todo.ID)
	if err != nil {
		t.Fatalf("UncompleteTodo() error = %v", err)
	}

	updated, err := GetTodoByID(db, todo.ID)
	if err != nil {
		t.Fatalf("GetTodoByID() error = %v", err)
	}

	if updated.IsComplete {
		t.Error("UncompleteTodo() did not mark todo as incomplete")
	}

	if updated.CompletedAt.Valid {
		t.Error("UncompleteTodo() did not clear completed_at timestamp")
	}
}

func TestListTodos(t *testing.T) {
	db := setupTestDB(t)

	tomorrow := time.Now().AddDate(0, 0, 1)

	_, err := CreateTodo(db, "Todo 1", []string{"work"}, nil)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	completedTodo, err := CreateTodo(db, "Completed todo", []string{}, nil)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	CompleteTodo(db, completedTodo.ID)

	_, err = CreateTodo(db, "Todo with due date", []string{"urgent"}, &tomorrow)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	t.Run("list all todos", func(t *testing.T) {
		todos, err := ListTodos(db, TodoListOptions{})
		if err != nil {
			t.Fatalf("ListTodos() error = %v", err)
		}

		if len(todos) != 3 {
			t.Errorf("ListTodos() returned %d todos, want 3", len(todos))
		}
	})

	t.Run("filter incomplete todos", func(t *testing.T) {
		todos, err := ListTodos(db, TodoListOptions{Incomplete: true})
		if err != nil {
			t.Fatalf("ListTodos() error = %v", err)
		}

		if len(todos) != 2 {
			t.Errorf("ListTodos() returned %d todos, want 2", len(todos))
		}

		for _, todo := range todos {
			if todo.IsComplete {
				t.Error("ListTodos() with Incomplete returned completed todo")
			}
		}
	})

	t.Run("filter by tag", func(t *testing.T) {
		todos, err := ListTodos(db, TodoListOptions{Tags: []string{"work"}})
		if err != nil {
			t.Fatalf("ListTodos() error = %v", err)
		}

		if len(todos) != 1 {
			t.Errorf("ListTodos() returned %d todos, want 1", len(todos))
		}
	})
}

func TestDeleteTodo(t *testing.T) {
	db := setupTestDB(t)

	todo, err := CreateTodo(db, "Todo to delete", []string{}, nil)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	t.Run("delete existing todo", func(t *testing.T) {
		err := DeleteTodo(db, todo.ID)
		if err != nil {
			t.Fatalf("DeleteTodo() error = %v", err)
		}

		_, err = GetTodoByID(db, todo.ID)
		if err == nil {
			t.Error("DeleteTodo() did not delete the todo")
		}
	})

	t.Run("delete non-existing todo", func(t *testing.T) {
		err := DeleteTodo(db, 99999)
		if err == nil {
			t.Error("DeleteTodo() expected error for non-existing todo, got nil")
		}
	})
}

func timePtr(t time.Time) *time.Time {
	return &t
}

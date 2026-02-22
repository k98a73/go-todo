package storage

import (
	"context"
	"os"
	"testing"

	"github.com/k98a73/go-todo/internal/domain"
)

func TestFileRepository_Create(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "todo*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte("[]"))
	tmpfile.Close()

	repo := NewFileRepository(tmpfile.Name())
	todo := &domain.Todo{Title: "Buy milk"}

	err = repo.Create(context.Background(), todo)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if todo.ID == 0 {
		t.Error("Expected ID to be assigned")
	}
}

func TestFileRepository_List(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "todo*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`[{"id":1,"title":"Test","completed":false}]`))
	tmpfile.Close()

	repo := NewFileRepository(tmpfile.Name())

	todos, err := repo.List(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}
}

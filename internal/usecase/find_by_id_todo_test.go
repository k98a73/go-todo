package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/k98a73/go-todo/internal/domain"
)

func TestFindByIDTodoUsecase_Execute(t *testing.T) {
	now := time.Now()
	mock := &MockRepository{
		todoList: []*domain.Todo{
			{ID: 1, Title: "Buy milk", Completed: false, CreatedAt: now, UpdatedAt: now},
			{ID: 2, Title: "Read book", Completed: true, CreatedAt: now, UpdatedAt: now},
		},
	}
	usecase := NewFindByIDTodoUsecase(mock)

	todo, err := usecase.Execute(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if todo.ID != 1 {
		t.Errorf("Expected ID 1, got %d", todo.ID)
	}
	if todo.Title != "Buy milk" {
		t.Errorf("Expected title 'Buy milk', got '%s'", todo.Title)
	}
}

func TestFindByIDTodoUsecase_Execute_NotFound(t *testing.T) {
	mock := &MockRepository{
		todoList: []*domain.Todo{},
	}
	usecase := NewFindByIDTodoUsecase(mock)

	_, err := usecase.Execute(context.Background(), 999)

	if err == nil {
		t.Error("Expected error for non-existent todo")
	}
}

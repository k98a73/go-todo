package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/k98a73/go-todo/internal/domain"
)

func TestListTodoUsecase_Execute(t *testing.T) {
	now := time.Now()
	mock := &MockRepository{
		todoList: []*domain.Todo{
			{ID: 1, Title: "Buy milk", Completed: false, CreatedAt: now, UpdatedAt: now},
			{ID: 2, Title: "Read book", Completed: true, CreatedAt: now, UpdatedAt: now},
		},
	}
	usecase := NewListTodoUsecase(mock)

	todoList, err := usecase.Execute(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(todoList) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todoList))
	}
	if todoList[0].Title != "Buy milk" {
		t.Errorf("Expected title 'Buy milk', got '%s'", todoList[0].Title)
	}
}

func TestListTodoUsecase_Execute_Empty(t *testing.T) {
	mock := &MockRepository{
		todoList: []*domain.Todo{},
	}
	usecase := NewListTodoUsecase(mock)

	todoList, err := usecase.Execute(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(todoList) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todoList))
	}
}

package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/k98a73/go-todo/internal/domain"
)

func TestUpdateTodoUsecase_Execute(t *testing.T) {
	now := time.Now()
	mock := &MockRepository{
		todoList: []*domain.Todo{
			{ID: 1, Title: "Buy milk", Completed: false, CreatedAt: now, UpdatedAt: now},
		},
	}
	usecase := NewUpdateTodoUsecase(mock)

	todo, err := usecase.Execute(context.Background(), 1, "Buy milk and eggs", true)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock.updateCalled {
		t.Error("Expected Update to be called")
	}
	if todo.Title != "Buy milk and eggs" {
		t.Errorf("Expected title 'Buy milk and eggs', got '%s'", todo.Title)
	}
	if !todo.Completed {
		t.Error("Expected completed to be true")
	}
	if todo.CreatedAt != now {
		t.Error("Expected CreatedAt to remain unchanged")
	}
	if !todo.UpdatedAt.After(now) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestUpdateTodoUsecase_Execute_EmptyTitle(t *testing.T) {
	now := time.Now()
	mock := &MockRepository{
		todoList: []*domain.Todo{
			{ID: 1, Title: "Buy milk", Completed: false, CreatedAt: now, UpdatedAt: now},
		},
	}
	usecase := NewUpdateTodoUsecase(mock)

	_, err := usecase.Execute(context.Background(), 1, "", false)

	if err == nil {
		t.Error("Expected error for empty title")
	}
}

func TestUpdateTodoUsecase_Execute_NotFound(t *testing.T) {
	mock := &MockRepository{
		todoList: []*domain.Todo{},
	}
	usecase := NewUpdateTodoUsecase(mock)

	_, err := usecase.Execute(context.Background(), 999, "Updated", false)

	if err == nil {
		t.Error("Expected error for non-existent todo")
	}
}

func TestUpdateTodoUsecase_Execute_RepoUpdateError(t *testing.T) {
	// Given: repo.Update がエラーを返すモック
	// When:  Execute を呼び出す
	// Then:  エラーが伝播する
	now := time.Now()
	mock := &MockRepository{
		todoList: []*domain.Todo{
			{ID: 1, Title: "Buy milk", Completed: false, CreatedAt: now, UpdatedAt: now},
		},
		updateErr: errors.New("storage failure"),
	}
	usecase := NewUpdateTodoUsecase(mock)

	_, err := usecase.Execute(context.Background(), 1, "Updated", true)

	if err == nil {
		t.Error("Expected error when repo.Update fails")
	}
}

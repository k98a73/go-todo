package usecase

import (
	"context"
	"testing"

	"github.com/k98a73/go-todo/internal/domain"
)

type MockRepository struct {
	createCalled bool
	createdTodo  *domain.Todo
	todoList     []*domain.Todo
}

func (m *MockRepository) Create(ctx context.Context, todo *domain.Todo) error {
	m.createCalled = true
	m.createdTodo = todo
	todo.ID = 1
	return nil
}

func (m *MockRepository) List(ctx context.Context) ([]*domain.Todo, error) {
	return m.todoList, nil
}

func (m *MockRepository) FindByID(ctx context.Context, id int) (*domain.Todo, error) {
	return nil, nil
}

func (m *MockRepository) Update(ctx context.Context, todo *domain.Todo) error {
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id int) error {
	return nil
}

// --- テスト ---
func TestCreateTodoUsecase_Execute(t *testing.T) {
	mock := &MockRepository{}
	usecase := NewCreateTodoUsecase(mock)

	todo, err := usecase.Execute(context.Background(), "Buy milk")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock.createCalled {
		t.Error("Expected Create to be called")
	}
	if todo.Title != "Buy milk" {
		t.Errorf("Expected title 'Buy milk', got '%s'", todo.Title)
	}
}

func TestCreateTodoUsecase_Execute_EmptyTitle(t *testing.T) {
	mock := &MockRepository{}
	usecase := NewCreateTodoUsecase(mock)

	_, err := usecase.Execute(context.Background(), "")

	if err == nil {
		t.Error("Expected error for empty title")
	}
}

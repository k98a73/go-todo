package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/k98a73/go-todo/internal/domain"
)

type mockCreateTodoUsecase struct {
	err error
}

func (m *mockCreateTodoUsecase) Execute(ctx context.Context, title string) (*domain.Todo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.Todo{
		ID:        1,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func TestCreateTodoHandler(t *testing.T) {
	mockUsecase := &mockCreateTodoUsecase{}
	handler := NewTodoHandler(mockUsecase)

	body := strings.NewReader(`{"title": "Buy milk"}`)
	req, _ := http.NewRequest("POST", "/todo", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.CreateTodo(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestCreateTodoHandler_EmptyTitle(t *testing.T) {
	mockUsecase := &mockCreateTodoUsecase{
		err: domain.ValidateTodo(&domain.Todo{Title: ""}),
	}
	handler := NewTodoHandler(mockUsecase)

	body := strings.NewReader(`{"title": ""}`)
	req, _ := http.NewRequest("POST", "/todo", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.CreateTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

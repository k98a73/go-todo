package http

import (
	"context"
	"fmt"
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
	handler := NewTodoHandler(mockUsecase, nil, nil, nil, nil)

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
	handler := NewTodoHandler(mockUsecase, nil, nil, nil, nil)

	body := strings.NewReader(`{"title": ""}`)
	req, _ := http.NewRequest("POST", "/todo", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.CreateTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

type mockListTodoUsecase struct {
	err   error
	todos []*domain.Todo
}

func (m *mockListTodoUsecase) Execute(ctx context.Context) ([]*domain.Todo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.todos, nil
}

func TestListTodoHandler(t *testing.T) {
	mockCreate := &mockCreateTodoUsecase{}
	mockList := &mockListTodoUsecase{
		todos: []*domain.Todo{
			{ID: 1, Title: "Buy milk"},
			{ID: 2, Title: "Go to gym"},
		},
	}
	// Note: NewTodoHandler will eventually need all usecases, but we'll update it incrementally
	handler := NewTodoHandler(mockCreate, mockList, nil, nil, nil)

	req, _ := http.NewRequest("GET", "/todo/list", nil)
	w := httptest.NewRecorder()

	handler.ListTodo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

type mockFindByIDTodoUsecase struct {
	err  error
	todo *domain.Todo
}

func (m *mockFindByIDTodoUsecase) Execute(ctx context.Context, id int) (*domain.Todo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.todo, nil
}

func TestFindByIDTodoHandler(t *testing.T) {
	mockFind := &mockFindByIDTodoUsecase{
		todo: &domain.Todo{ID: 1, Title: "Buy milk"},
	}
	handler := NewTodoHandler(nil, nil, mockFind, nil, nil)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	req.SetPathValue("id", "1") // Simulate routing path value
	w := httptest.NewRecorder()

	handler.FindByIDTodo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestFindByIDTodoHandler_InvalidID(t *testing.T) {
	handler := NewTodoHandler(nil, nil, nil, nil, nil)

	req, _ := http.NewRequest("GET", "/todo/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	handler.FindByIDTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}
}

type mockUpdateTodoUsecase struct {
	err  error
	todo *domain.Todo
}

func (m *mockUpdateTodoUsecase) Execute(ctx context.Context, id int, title string, completed bool) (*domain.Todo, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.todo != nil {
		m.todo.Title = title
		m.todo.Completed = completed
		return m.todo, nil
	}
	return &domain.Todo{
		ID:        id,
		Title:     title,
		Completed: completed,
		UpdatedAt: time.Now(),
	}, nil
}

func TestUpdateTodoHandler(t *testing.T) {
	mockUpdate := &mockUpdateTodoUsecase{}
	handler := NewTodoHandler(nil, nil, nil, mockUpdate, nil)

	body := strings.NewReader(`{"title": "Updated title", "completed": true}`)
	req, _ := http.NewRequest("PUT", "/todo/1", body)
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.UpdateTodo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateTodoHandler_InvalidBody(t *testing.T) {
	handler := NewTodoHandler(nil, nil, nil, nil, nil)
	body := strings.NewReader(`invalid json`)
	req, _ := http.NewRequest("PUT", "/todo/1", body)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	handler.UpdateTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

type mockDeleteTodoUsecase struct {
	err error
}

func (m *mockDeleteTodoUsecase) Execute(ctx context.Context, id int) error {
	return m.err
}

func TestDeleteTodoHandler(t *testing.T) {
	mockDelete := &mockDeleteTodoUsecase{}
	handler := NewTodoHandler(nil, nil, nil, nil, mockDelete)

	req, _ := http.NewRequest("DELETE", "/todo/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteTodo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteTodoHandler_InvalidID(t *testing.T) {
	handler := NewTodoHandler(nil, nil, nil, nil, nil)

	req, _ := http.NewRequest("DELETE", "/todo/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	handler.DeleteTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}
}

func TestCreateTodoHandler_InvalidJSON(t *testing.T) {
	// Given: 不正なリクエストボディ
	// When:  CreateTodo を呼び出す
	// Then:  400 Bad Request が返る
	handler := NewTodoHandler(nil, nil, nil, nil, nil)

	body := strings.NewReader(`not json`)
	req, _ := http.NewRequest("POST", "/todo", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateTodoHandler_UsecaseError(t *testing.T) {
	// Given: usecase が内部エラーを返すモック
	// When:  CreateTodo を呼び出す
	// Then:  500 Internal Server Error が返る
	mockUsecase := &mockCreateTodoUsecase{err: fmt.Errorf("repository error")}
	handler := NewTodoHandler(mockUsecase, nil, nil, nil, nil)

	body := strings.NewReader(`{"title": "Some title"}`)
	req, _ := http.NewRequest("POST", "/todo", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTodo(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestListTodoHandler_UsecaseError(t *testing.T) {
	// Given: usecase が内部エラーを返すモック
	// When:  ListTodo を呼び出す
	// Then:  500 Internal Server Error が返る
	mockList := &mockListTodoUsecase{err: fmt.Errorf("repository error")}
	handler := NewTodoHandler(nil, mockList, nil, nil, nil)

	req, _ := http.NewRequest("GET", "/todo/list", nil)
	w := httptest.NewRecorder()

	handler.ListTodo(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestFindByIDTodoHandler_NotFound(t *testing.T) {
	// Given: usecase が "todo not found" エラーを返すモック
	// When:  FindByIDTodo を呼び出す
	// Then:  404 Not Found が返る
	mockFind := &mockFindByIDTodoUsecase{err: fmt.Errorf("todo not found")}
	handler := NewTodoHandler(nil, nil, mockFind, nil, nil)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.FindByIDTodo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestFindByIDTodoHandler_UsecaseError(t *testing.T) {
	// Given: usecase が内部エラーを返すモック
	// When:  FindByIDTodo を呼び出す
	// Then:  500 Internal Server Error が返る
	mockFind := &mockFindByIDTodoUsecase{err: fmt.Errorf("internal error")}
	handler := NewTodoHandler(nil, nil, mockFind, nil, nil)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.FindByIDTodo(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestFindByIDTodoHandler_NilTodo(t *testing.T) {
	// Given: usecase が nil Todo を返すモック
	// When:  FindByIDTodo を呼び出す
	// Then:  404 Not Found が返る
	mockFind := &mockFindByIDTodoUsecase{todo: nil}
	handler := NewTodoHandler(nil, nil, mockFind, nil, nil)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.FindByIDTodo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for nil todo, got %d", w.Code)
	}
}

func TestUpdateTodoHandler_InvalidID(t *testing.T) {
	// Given: 不正なID
	// When:  UpdateTodo を呼び出す
	// Then:  400 Bad Request が返る
	handler := NewTodoHandler(nil, nil, nil, nil, nil)

	body := strings.NewReader(`{"title": "test", "completed": false}`)
	req, _ := http.NewRequest("PUT", "/todo/abc", body)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTodoHandler_NotFound(t *testing.T) {
	// Given: usecase が "todo not found" エラーを返すモック
	// When:  UpdateTodo を呼び出す
	// Then:  404 Not Found が返る
	mockUpdate := &mockUpdateTodoUsecase{err: fmt.Errorf("todo not found")}
	handler := NewTodoHandler(nil, nil, nil, mockUpdate, nil)

	body := strings.NewReader(`{"title": "test", "completed": false}`)
	req, _ := http.NewRequest("PUT", "/todo/1", body)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateTodoHandler_TitleEmpty(t *testing.T) {
	// Given: usecase が "title cannot be empty" を返すモック
	// When:  UpdateTodo を呼び出す
	// Then:  400 Bad Request が返る
	mockUpdate := &mockUpdateTodoUsecase{err: fmt.Errorf("title cannot be empty")}
	handler := NewTodoHandler(nil, nil, nil, mockUpdate, nil)

	body := strings.NewReader(`{"title": "", "completed": false}`)
	req, _ := http.NewRequest("PUT", "/todo/1", body)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTodoHandler_TitleTooLong(t *testing.T) {
	// Given: usecase が "title too long" を返すモック
	// When:  UpdateTodo を呼び出す
	// Then:  400 Bad Request が返る
	mockUpdate := &mockUpdateTodoUsecase{err: fmt.Errorf("title too long")}
	handler := NewTodoHandler(nil, nil, nil, mockUpdate, nil)

	body := strings.NewReader(`{"title": "long", "completed": false}`)
	req, _ := http.NewRequest("PUT", "/todo/1", body)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTodoHandler_UsecaseError(t *testing.T) {
	// Given: usecase が内部エラーを返すモック
	// When:  UpdateTodo を呼び出す
	// Then:  500 Internal Server Error が返る
	mockUpdate := &mockUpdateTodoUsecase{err: fmt.Errorf("internal error")}
	handler := NewTodoHandler(nil, nil, nil, mockUpdate, nil)

	body := strings.NewReader(`{"title": "test", "completed": false}`)
	req, _ := http.NewRequest("PUT", "/todo/1", body)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestDeleteTodoHandler_NotFound(t *testing.T) {
	// Given: usecase が "todo not found" エラーを返すモック
	// When:  DeleteTodo を呼び出す
	// Then:  404 Not Found が返る
	mockDelete := &mockDeleteTodoUsecase{err: fmt.Errorf("todo not found")}
	handler := NewTodoHandler(nil, nil, nil, nil, mockDelete)

	req, _ := http.NewRequest("DELETE", "/todo/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteTodo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteTodoHandler_UsecaseError(t *testing.T) {
	// Given: usecase が内部エラーを返すモック
	// When:  DeleteTodo を呼び出す
	// Then:  500 Internal Server Error が返る
	mockDelete := &mockDeleteTodoUsecase{err: fmt.Errorf("internal error")}
	handler := NewTodoHandler(nil, nil, nil, nil, mockDelete)

	req, _ := http.NewRequest("DELETE", "/todo/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.DeleteTodo(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

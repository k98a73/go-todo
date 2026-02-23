package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k98a73/go-todo/internal/domain"
)

type CreateTodoUsecase interface {
	Execute(ctx context.Context, title string) (*domain.Todo, error)
}

type ListTodoUsecase interface {
	Execute(ctx context.Context) ([]*domain.Todo, error)
}

type FindByIDTodoUsecase interface {
	Execute(ctx context.Context, id int) (*domain.Todo, error)
}

type UpdateTodoUsecase interface {
	Execute(ctx context.Context, id int, title string, completed bool) (*domain.Todo, error)
}

type DeleteTodoUsecase interface {
	Execute(ctx context.Context, id int) error
}

type TodoHandler struct {
	createUsecase   CreateTodoUsecase
	listUsecase     ListTodoUsecase
	findByIDUsecase FindByIDTodoUsecase
	updateUsecase   UpdateTodoUsecase
	deleteUsecase   DeleteTodoUsecase
}

func NewTodoHandler(create CreateTodoUsecase, list ListTodoUsecase, findByID FindByIDTodoUsecase, update UpdateTodoUsecase, del DeleteTodoUsecase) *TodoHandler {
	return &TodoHandler{
		createUsecase:   create,
		listUsecase:     list,
		findByIDUsecase: findByID,
		updateUsecase:   update,
		deleteUsecase:   del,
	}
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := h.createUsecase.Execute(r.Context(), req.Title)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) ListTodo(w http.ResponseWriter, r *http.Request) {
	todos, err := h.listUsecase.Execute(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) FindByIDTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := h.findByIDUsecase.Execute(r.Context(), id)
	if err != nil {
		if err.Error() == "todo not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

type UpdateTodoRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo, err := h.updateUsecase.Execute(r.Context(), id, req.Title, req.Completed)
	if err != nil {
		if err.Error() == "todo not found" {
			w.WriteHeader(http.StatusNotFound)
		} else if err.Error() == "title cannot be empty" || err.Error() == "title too long" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.deleteUsecase.Execute(r.Context(), id); err != nil {
		if err.Error() == "todo not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "todo deleted successfully"})
}

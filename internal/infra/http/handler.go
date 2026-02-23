package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/k98a73/go-todo/internal/domain"
)

type CreateTodoUsecase interface {
	Execute(ctx context.Context, title string) (*domain.Todo, error)
}

type TodoHandler struct {
	createUsecase CreateTodoUsecase
}

func NewTodoHandler(u CreateTodoUsecase) *TodoHandler {
	return &TodoHandler{createUsecase: u}
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

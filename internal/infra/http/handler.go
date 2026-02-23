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
	createTodoUsecase CreateTodoUsecase
}

func NewTodoHandler(createUsecase CreateTodoUsecase) *TodoHandler {
	return &TodoHandler{
		createTodoUsecase: createUsecase,
	}
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	todo, err := h.createTodoUsecase.Execute(r.Context(), req.Title)
	if err != nil {
		if err.Error() == "title cannot be empty" || err.Error() == "title too long" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

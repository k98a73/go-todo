package main

import (
	"log"
	"net/http"

	http_infra "github.com/k98a73/go-todo/internal/infra/http"
	"github.com/k98a73/go-todo/internal/infra/storage"
	"github.com/k98a73/go-todo/internal/usecase"
)

func main() {
	repo := storage.NewFileRepository("todos.json")
	createUsecase := usecase.NewCreateTodoUsecase(repo)
	todoHandler := http_infra.NewTodoHandler(createUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /todo", todoHandler.CreateTodo)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

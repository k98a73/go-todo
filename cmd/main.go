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
	listUsecase := usecase.NewListTodoUsecase(repo)
	findByIDUsecase := usecase.NewFindByIDTodoUsecase(repo)
	updateUsecase := usecase.NewUpdateTodoUsecase(repo)
	deleteUsecase := usecase.NewDeleteTodoUsecase(repo)
	todoHandler := http_infra.NewTodoHandler(createUsecase, listUsecase, findByIDUsecase, updateUsecase, deleteUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /todo", todoHandler.CreateTodo)
	mux.HandleFunc("GET /todo/list", todoHandler.ListTodo)
	mux.HandleFunc("GET /todo/{id}", todoHandler.FindByIDTodo)
	mux.HandleFunc("PUT /todo/{id}", todoHandler.UpdateTodo)
	mux.HandleFunc("DELETE /todo/{id}", todoHandler.DeleteTodo)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

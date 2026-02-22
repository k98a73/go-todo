package usecase

import (
	"context"
	"time"

	"github.com/k98a73/go-todo/internal/domain"
)

type CreateTodoUsecase struct {
	repo domain.IRepository
}

func NewCreateTodoUsecase(repo domain.IRepository) *CreateTodoUsecase {
	return &CreateTodoUsecase{repo: repo}
}

func (u *CreateTodoUsecase) Execute(ctx context.Context, title string) (*domain.Todo, error) {
	now := time.Now()
	todo := &domain.Todo{
		Title:     title,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := domain.ValidateTodo(todo); err != nil {
		return nil, err
	}

	if err := u.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

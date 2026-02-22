package usecase

import (
	"context"

	"github.com/k98a73/go-todo/internal/domain"
)

type ListTodoUsecase struct {
	repo domain.IRepository
}

func NewListTodoUsecase(repo domain.IRepository) *ListTodoUsecase {
	return &ListTodoUsecase{repo: repo}
}

func (u *ListTodoUsecase) Execute(ctx context.Context) ([]*domain.Todo, error) {
	return u.repo.List(ctx)
}

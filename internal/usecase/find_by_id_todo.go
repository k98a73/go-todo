package usecase

import (
	"context"

	"github.com/k98a73/go-todo/internal/domain"
)

type FindByIDTodoUsecase struct {
	repo domain.IRepository
}

func NewFindByIDTodoUsecase(repo domain.IRepository) *FindByIDTodoUsecase {
	return &FindByIDTodoUsecase{repo: repo}
}

func (u *FindByIDTodoUsecase) Execute(ctx context.Context, id int) (*domain.Todo, error) {
	return u.repo.FindByID(ctx, id)
}

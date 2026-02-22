package usecase

import (
	"context"

	"github.com/k98a73/go-todo/internal/domain"
)

type DeleteTodoUsecase struct {
	repo domain.IRepository
}

func NewDeleteTodoUsecase(repo domain.IRepository) *DeleteTodoUsecase {
	return &DeleteTodoUsecase{repo: repo}
}

func (u *DeleteTodoUsecase) Execute(ctx context.Context, id int) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

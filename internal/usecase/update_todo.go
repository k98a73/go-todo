package usecase

import (
	"context"
	"time"

	"github.com/k98a73/go-todo/internal/domain"
)

type UpdateTodoUsecase struct {
	repo domain.IRepository
}

func NewUpdateTodoUsecase(repo domain.IRepository) *UpdateTodoUsecase {
	return &UpdateTodoUsecase{repo: repo}
}

func (u *UpdateTodoUsecase) Execute(ctx context.Context, id int, title string, completed bool) (*domain.Todo, error) {
	todo, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo.Title = title
	todo.Completed = completed
	todo.UpdatedAt = time.Now()

	if err := domain.ValidateTodo(todo); err != nil {
		return nil, err
	}

	if err := u.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

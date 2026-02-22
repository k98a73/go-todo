package domain

import (
	"context"
	"errors"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ValidateTodo(t *Todo) error {
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(t.Title) > 255 {
		return errors.New("title too long")
	}
	return nil
}

type IRepository interface {
	Create(ctx context.Context, todo *Todo) error
	List(ctx context.Context) ([]*Todo, error)
	FindByID(ctx context.Context, id int) (*Todo, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id int) error
}

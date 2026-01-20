package domain

import (
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

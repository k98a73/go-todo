package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/k98a73/go-todo/internal/domain"
)

type FileRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewFileRepository(filePath string) *FileRepository {
	return &FileRepository{
		filePath: filePath,
	}
}

func (r *FileRepository) load() ([]*domain.Todo, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*domain.Todo{}, nil
		}
		return nil, err
	}

	var todos []*domain.Todo
	if len(data) == 0 {
		return []*domain.Todo{}, nil
	}

	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *FileRepository) save(todos []*domain.Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func (r *FileRepository) Create(ctx context.Context, todo *domain.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todos, err := r.load()
	if err != nil {
		return err
	}

	maxID := 0
	for _, t := range todos {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	todo.ID = maxID + 1

	todos = append(todos, todo)

	return r.save(todos)
}

func (r *FileRepository) List(ctx context.Context) ([]*domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.load()
}

func (r *FileRepository) FindByID(ctx context.Context, id int) (*domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos, err := r.load()
	if err != nil {
		return nil, err
	}

	for _, t := range todos {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, errors.New("todo not found")
}

func (r *FileRepository) Update(ctx context.Context, todo *domain.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todos, err := r.load()
	if err != nil {
		return err
	}

	for i, t := range todos {
		if t.ID == todo.ID {
			todos[i] = todo
			return r.save(todos)
		}
	}

	return errors.New("todo not found")
}

func (r *FileRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todos, err := r.load()
	if err != nil {
		return err
	}

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return r.save(todos)
		}
	}

	return errors.New("todo not found")
}

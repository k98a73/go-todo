package usecase

import (
	"context"
	"testing"
)

func TestDeleteTodoUsecase_Execute(t *testing.T) {
	mock := &MockRepository{}
	usecase := NewDeleteTodoUsecase(mock)

	err := usecase.Execute(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mock.deleteCalled {
		t.Error("Expected Delete to be called")
	}
	if mock.deletedID != 1 {
		t.Errorf("Expected deleted ID 1, got %d", mock.deletedID)
	}
}

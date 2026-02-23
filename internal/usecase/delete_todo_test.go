package usecase

import (
	"context"
	"errors"
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

func TestDeleteTodoUsecase_Execute_RepoError(t *testing.T) {
	// Given: repo.Delete がエラーを返すモック
	// When:  Execute を呼び出す
	// Then:  エラーが伝播する
	mock := &MockRepository{deleteErr: errors.New("storage failure")}
	usecase := NewDeleteTodoUsecase(mock)

	err := usecase.Execute(context.Background(), 1)

	if err == nil {
		t.Error("Expected error when repo.Delete fails")
	}
}

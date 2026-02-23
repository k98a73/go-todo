package storage

import (
	"context"
	"os"
	"testing"

	"github.com/k98a73/go-todo/internal/domain"
)

func newTempRepo(t *testing.T, content string) (*FileRepository, func()) {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "todo*.json")
	if err != nil {
		t.Fatal(err)
	}
	if content != "" {
		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	tmpfile.Close()
	return NewFileRepository(tmpfile.Name()), func() { os.Remove(tmpfile.Name()) }
}

func TestFileRepository_Create(t *testing.T) {
	// Given: 空のリポジトリ
	// When:  Create を呼び出す
	// Then:  エラーなし・IDが割り当てられる
	repo, cleanup := newTempRepo(t, "[]")
	defer cleanup()

	todo := &domain.Todo{Title: "Buy milk"}
	err := repo.Create(context.Background(), todo)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if todo.ID == 0 {
		t.Error("Expected ID to be assigned")
	}
}

func TestFileRepository_Create_AutoIncrement(t *testing.T) {
	// Given: 既存Todoが1件あるリポジトリ
	// When:  さらに Create を呼び出す
	// Then:  IDが maxID+1 になる
	repo, cleanup := newTempRepo(t, `[{"id":5,"title":"Existing","completed":false}]`)
	defer cleanup()

	todo := &domain.Todo{Title: "New todo"}
	err := repo.Create(context.Background(), todo)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if todo.ID != 6 {
		t.Errorf("Expected ID 6, got %d", todo.ID)
	}
}

func TestFileRepository_Create_LoadError(t *testing.T) {
	// Given: 不正なJSONが書かれたファイル
	// When:  Create を呼び出す
	// Then:  load エラーが伝播する
	repo, cleanup := newTempRepo(t, "invalid json")
	defer cleanup()

	todo := &domain.Todo{Title: "Test"}
	err := repo.Create(context.Background(), todo)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestFileRepository_List(t *testing.T) {
	// Given: 1件のTodoが入ったファイル
	// When:  List を呼び出す
	// Then:  1件のスライスが返る
	repo, cleanup := newTempRepo(t, `[{"id":1,"title":"Test","completed":false}]`)
	defer cleanup()

	todos, err := repo.List(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}
}

func TestFileRepository_List_FileNotExist(t *testing.T) {
	// Given: 存在しないファイルパスのリポジトリ
	// When:  List を呼び出す
	// Then:  エラーなし・空スライスが返る
	repo := NewFileRepository("/tmp/nonexistent_todo_file_12345.json")

	todos, err := repo.List(context.Background())

	if err != nil {
		t.Errorf("Expected no error for non-existent file, got %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}
}

func TestFileRepository_List_EmptyFile(t *testing.T) {
	// Given: 空ファイル（0バイト）のリポジトリ
	// When:  List を呼び出す
	// Then:  エラーなし・空スライスが返る
	repo, cleanup := newTempRepo(t, "")
	defer cleanup()

	todos, err := repo.List(context.Background())

	if err != nil {
		t.Errorf("Expected no error for empty file, got %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}
}

func TestFileRepository_List_InvalidJSON(t *testing.T) {
	// Given: 不正なJSONが書かれたファイル
	// When:  List を呼び出す
	// Then:  json.Unmarshal エラーが返る
	repo, cleanup := newTempRepo(t, "not-valid-json")
	defer cleanup()

	_, err := repo.List(context.Background())

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestFileRepository_FindByID(t *testing.T) {
	// Given: 2件のTodoが入ったリポジトリ
	// When:  存在するIDで FindByID を呼び出す
	// Then:  該当のTodoが返る
	repo, cleanup := newTempRepo(t, `[{"id":1,"title":"Buy milk","completed":false},{"id":2,"title":"Read book","completed":true}]`)
	defer cleanup()

	todo, err := repo.FindByID(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if todo == nil {
		t.Fatal("Expected todo, got nil")
	}
	if todo.ID != 1 {
		t.Errorf("Expected ID 1, got %d", todo.ID)
	}
	if todo.Title != "Buy milk" {
		t.Errorf("Expected title 'Buy milk', got '%s'", todo.Title)
	}
}

func TestFileRepository_FindByID_NotFound(t *testing.T) {
	// Given: 1件のTodoが入ったリポジトリ
	// When:  存在しないIDで FindByID を呼び出す
	// Then:  "todo not found" エラーが返る
	repo, cleanup := newTempRepo(t, `[{"id":1,"title":"Buy milk","completed":false}]`)
	defer cleanup()

	_, err := repo.FindByID(context.Background(), 999)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "todo not found" {
		t.Errorf("Expected 'todo not found', got '%s'", err.Error())
	}
}

func TestFileRepository_FindByID_LoadError(t *testing.T) {
	// Given: 不正なJSONが書かれたファイル
	// When:  FindByID を呼び出す
	// Then:  load エラーが伝播する
	repo, cleanup := newTempRepo(t, "bad json")
	defer cleanup()

	_, err := repo.FindByID(context.Background(), 1)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestFileRepository_Update(t *testing.T) {
	// Given: 1件のTodoが入ったリポジトリ
	// When:  Update を呼び出す
	// Then:  エラーなし・内容が更新される
	repo, cleanup := newTempRepo(t, `[{"id":1,"title":"Buy milk","completed":false}]`)
	defer cleanup()

	updated := &domain.Todo{ID: 1, Title: "Buy milk and eggs", Completed: true}
	err := repo.Update(context.Background(), updated)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	todo, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if todo.Title != "Buy milk and eggs" {
		t.Errorf("Expected updated title, got '%s'", todo.Title)
	}
	if !todo.Completed {
		t.Error("Expected completed to be true")
	}
}

func TestFileRepository_Update_NotFound(t *testing.T) {
	// Given: 1件のTodoが入ったリポジトリ
	// When:  存在しないIDで Update を呼び出す
	// Then:  "todo not found" エラーが返る
	repo, cleanup := newTempRepo(t, `[{"id":1,"title":"Buy milk","completed":false}]`)
	defer cleanup()

	err := repo.Update(context.Background(), &domain.Todo{ID: 999, Title: "Ghost"})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "todo not found" {
		t.Errorf("Expected 'todo not found', got '%s'", err.Error())
	}
}

func TestFileRepository_Update_LoadError(t *testing.T) {
	// Given: 不正なJSONが書かれたファイル
	// When:  Update を呼び出す
	// Then:  load エラーが伝播する
	repo, cleanup := newTempRepo(t, "bad json")
	defer cleanup()

	err := repo.Update(context.Background(), &domain.Todo{ID: 1, Title: "Test"})

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestFileRepository_Delete(t *testing.T) {
	// Given: 2件のTodoが入ったリポジトリ
	// When:  Delete を呼び出す
	// Then:  エラーなし・対象が削除される
	repo, cleanup := newTempRepo(t, `[{"id":1,"title":"Buy milk","completed":false},{"id":2,"title":"Read book","completed":true}]`)
	defer cleanup()

	err := repo.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	todos, _ := repo.List(context.Background())
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo after delete, got %d", len(todos))
	}
	if todos[0].ID != 2 {
		t.Errorf("Expected remaining todo ID 2, got %d", todos[0].ID)
	}
}

func TestFileRepository_Delete_NotFound(t *testing.T) {
	// Given: 1件のTodoが入ったリポジトリ
	// When:  存在しないIDで Delete を呼び出す
	// Then:  "todo not found" エラーが返る
	repo, cleanup := newTempRepo(t, `[{"id":1,"title":"Buy milk","completed":false}]`)
	defer cleanup()

	err := repo.Delete(context.Background(), 999)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "todo not found" {
		t.Errorf("Expected 'todo not found', got '%s'", err.Error())
	}
}

func TestFileRepository_Delete_LoadError(t *testing.T) {
	// Given: 不正なJSONが書かれたファイル
	// When:  Delete を呼び出す
	// Then:  load エラーが伝播する
	repo, cleanup := newTempRepo(t, "bad json")
	defer cleanup()

	err := repo.Delete(context.Background(), 1)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

# テスト戦略

Goプロジェクトのテスト方法とテスト設計についてまとめます。

## テスト戦略概要

### テストピラミッド

```
        統合テスト（E2E）
       ↑
    API テスト
    ↑
   ユニットテスト
   ↑
```

### 各テストレベルの位置付け

| レベル | 範囲 | 実装 | 実行速度 |
|--------|------|------|--------|
| **ユニットテスト** | 単一関数・メソッド | Go `testing` パッケージ | 高速 |
| **API テスト** | HTTPエンドポイント | Go `testing` + `httptest` | 中速 |
| **統合テスト** | 複数コンポーネント | 実際のファイルI/O、DB | 低速 |

## ユニットテスト

### テストファイルの配置

```
internal/
├── domain/
│   ├── entity.go
│   └── entity_test.go      # テストはソースファイルと同じパッケージ
├── usecase/
│   ├── create_todo.go
│   └── create_todo_test.go
└── infra/
    ├── storage/
    │   ├── file_storage.go
    │   └── file_storage_test.go
```

### テストの基本構造

```go
package domain

import "testing"

// テスト関数の命名: Test + テスト対象
func TestValidateTodoTitle(t *testing.T) {
    // Arrange: テスト準備
    title := "Buy milk"
    
    // Act: テスト実行
    err := ValidateTodoTitle(title)
    
    // Assert: 検証
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

### テーブル駆動テスト（Table-Driven Tests）

複数のケースを効率的にテスト：

```go
func TestValidateTodoTitle(t *testing.T) {
    tests := []struct {
        name    string
        title   string
        wantErr bool
    }{
        {
            name:    "valid title",
            title:   "Buy milk",
            wantErr: false,
        },
        {
            name:    "empty title",
            title:   "",
            wantErr: true,
        },
        {
            name:    "very long title",
            title:   strings.Repeat("a", 300),
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateTodoTitle(tt.title)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateTodoTitle() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### モック・スタブの作成

インターフェースを使用したモック例：

```go
// モックリポジトリ
type MockTodoRepository struct {
    todoList []*Todo
}

func (m *MockTodoRepository) FindByID(ctx context.Context, id int) (*Todo, error) {
    for _, t := range m.todoList {
        if t.ID == id {
            return t, nil
        }
    }
    return nil, errors.New("not found")
}

// テストで使用
func TestFindByID(t *testing.T) {
    mock := &MockTodoRepository{
        todoList: []*Todo{
            {ID: 1, Title: "Test"},
        },
    }
    
    todo, err := mock.FindByID(context.Background(), 1)
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if todo.Title != "Test" {
        t.Errorf("Expected 'Test', got '%s'", todo.Title)
    }
}
```

## API テスト

### httptest を使用したテスト

```go
package http

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateTodo(t *testing.T) {
    // ハンドラーの初期化
    handler := NewTodoHandler(mockUsecase)
    
    // リクエストの作成
    body := strings.NewReader(`{"title": "Buy milk"}`)
    req, _ := http.NewRequest("POST", "/todo", body)
    req.Header.Set("Content-Type", "application/json")
    
    // レスポンスの記録
    w := httptest.NewRecorder()
    handler.CreateTodo(w, req)
    
    // 検証
    if w.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", w.Code)
    }
}
```

## 統合テスト

### ファイルI/O を含むテスト

```go
package infra

import (
    "os"
    "testing"
    "tempfile"
)

func TestFileStoragePersistence(t *testing.T) {
    // 一時ファイルの作成
    tmpfile, err := tempfile.CreateTemp("", "todo")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpfile.Name())
    
    // 保存
    storage := NewFileStorage(tmpfile.Name())
    todo := &Todo{Title: "Test"}
    if err := storage.Create(context.Background(), todo); err != nil {
        t.Fatalf("Create failed: %v", err)
    }
    
    // 読み込みで検証
    todoList, _ := storage.List(context.Background())
    if len(todoList) != 1 {
        t.Errorf("Expected 1 todo, got %d", len(todoList))
    }
}
```

## テストコマンド

### 全テストの実行

```bash
go test ./...
```

### 特定パッケージのテスト

```bash
go test ./internal/domain
```

### 詳細出力

```bash
go test -v ./...
```

### カバレッジ測定

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 特定テストのみ実行

```bash
go test -run TestValidateTodoTitle ./...
```

## 手動テスト（curl）

開発中の動作確認に使用：

### TODO 作成

```bash
curl -X POST http://localhost:8080/todo \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy milk"}'
```

### TODO 一覧

```bash
curl http://localhost:8080/todo/list
```

### TODO 更新

```bash
curl -X PUT http://localhost:8080/todo/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy milk", "completed": true}'
```

### TODO 削除

```bash
curl -X DELETE http://localhost:8080/todo/1
```

## ベストプラクティス

1. **テストごとに Arrange → Act → Assert を明確に**
2. **テーブル駆動テストで複数ケースをカバー**
3. **モックを使ってテスト対象を分離**
4. **エラーケースもテスト（正常系+異常系）**
5. **テスト名は目的を明確に**
6. **テストから実装の品質を保証**

## 参考リンク

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/CodeReviewComments#tests)

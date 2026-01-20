# 実装計画

Go Todoアプリケーション開発の実装順序と各段階での学習ポイントをまとめます。

**本計画はTest-Driven Development（TDD）を前提としています。**

## TDD の鉄則

```
失敗するテストなしに本番コードを書かない
```

すべての機能実装は **Red-Green-Refactor** サイクルに従います：

1. **RED** - 失敗するテストを書く
2. **Verify RED** - テストが正しく失敗することを確認
3. **GREEN** - テストを通す最小限のコードを書く
4. **Verify GREEN** - すべてのテストが通ることを確認
5. **REFACTOR** - コードを整理（テストは常にグリーンを維持）

**参考資料:** [TDD スキル](../../.agent/skills/test-driven-development/SKILL.md)

---

## 全体フロー

```
段階1: 基礎準備
  ↓
段階2: ドメイン層実装（TDD）
  ↓
段階3: ユースケース層実装（TDD）
  ↓
段階4: インフラストラクチャ層実装（TDD）
  ↓
段階5: HTTP層実装（TDD）
  ↓
段階6: 統合テスト・最適化
```

---

## 段階1: 基礎準備（学習）

### 1.1 開発環境セットアップ

**TODO:**
- [x] Go のインストール
- [x] リポジトリの作成
- [x] リポジトリのクローン
- [x] プロジェクトディレクトリ初期化
- [x] `go.mod` ファイル生成
- [x] mainブランチの作成

**学習ポイント:**
- Go のパッケージシステムの理解
- Module の基本
- Jujutsu（jj）によるバージョン管理

### 1.2 Go 言語の基礎学習

**TODO:**
- [x] struct（構造体）を学習
- [x] interface（インターフェース）を学習
- [x] time.Time を学習
- [x] context.Context を学習
- [x] error インターフェースを学習
- [x] defer を学習
- [x] ポインタを学習

**参考資料:** [GO_BASICS.md](./GO_BASICS.md)

**実装コマンド:**
詳細なコマンドリファレンスについては、[COMMANDS.md](./COMMANDS.md#初期セットアップコマンド) を参照してください。

---

## 段階2: ドメイン層実装（TDD）

### 2.1 バリデーション関数のテストと実装

#### 2.1.1 RED: 失敗するテストを書く

**TODO:**
- [x] `internal/domain/entity_test.go` を作成
- [x] `ValidateTodo` 関数のテストを書く（関数はまだ存在しない）

**テストを先に書く:**
```go
package domain

import (
    "strings"
    "testing"
)

func TestValidateTodo(t *testing.T) {
    tests := []struct {
        name    string
        title   string
        wantErr bool
        errMsg  string
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
            errMsg:  "title cannot be empty",
        },
        {
            name:    "title too long",
            title:   strings.Repeat("a", 256),
            wantErr: true,
            errMsg:  "title too long",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            todo := &Todo{Title: tt.title}
            err := ValidateTodo(todo)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateTodo() error = %v, wantErr %v", err, tt.wantErr)
            }
            if tt.wantErr && err != nil && err.Error() != tt.errMsg {
                t.Errorf("ValidateTodo() error = %v, want %v", err.Error(), tt.errMsg)
            }
        })
    }
}
```

#### 2.1.2 Verify RED: テストが失敗することを確認

**TODO:**
- [x] テストが失敗することを確認

```bash
go test ./internal/domain -v
# コンパイルエラー: Todo, ValidateTodo が未定義
```

#### 2.1.3 GREEN: テストを通す最小限のコードを書く

**TODO:**
- [x] `internal/domain/entity.go` を作成
- [x] `Todo` 構造体を定義
- [x] `ValidateTodo` 関数を実装

**実装:**
```go
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
```

#### 2.1.4 Verify GREEN: すべてのテストが通ることを確認

```bash
go test ./internal/domain -v
# PASS
```

**学習ポイント:**
- テーブル駆動テスト
- エラー返却パターン
- struct フィールド定義
- struct タグ（JSON マーシャリング）

### 2.2 リポジトリインターフェース定義

**TODO:**
- [ ] `IRepository` インターフェース定義
- [ ] 5つのメソッド（Create, List, FindByID, Update, Delete）

**期待される実装:**
```go
type IRepository interface {
    Create(ctx context.Context, todo *Todo) error
    List(ctx context.Context) ([]*Todo, error)
    FindByID(ctx context.Context, id int) (*Todo, error)
    Update(ctx context.Context, todo *Todo) error
    Delete(ctx context.Context, id int) error
}
```

**学習ポイント:**
- interface の定義
- context.Context の導入
- メソッドシグネチャ設計

---

## 段階3: ユースケース層実装（TDD）

### 3.1 Create ユースケース

#### 3.1.1 RED: 失敗するテストを書く

**TODO:**
- [ ] `internal/usecase/create_todo_test.go` を作成
- [ ] モックリポジトリを定義
- [ ] `CreateTodoUsecase` のテストを書く

**テストを先に書く:**
```go
package usecase

import (
    "context"
    "testing"

    "github.com/k98a73/go-todo/internal/domain"
)

type MockRepository struct {
    createCalled bool
    createdTodo  *domain.Todo
}

func (m *MockRepository) Create(ctx context.Context, todo *domain.Todo) error {
    m.createCalled = true
    m.createdTodo = todo
    todo.ID = 1  // ID を割り当てる
    return nil
}

// 他のメソッドも空実装...

func TestCreateTodoUsecase_Execute(t *testing.T) {
    mock := &MockRepository{}
    usecase := NewCreateTodoUsecase(mock)

    todo, err := usecase.Execute(context.Background(), "Buy milk")

    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if !mock.createCalled {
        t.Error("Expected Create to be called")
    }
    if todo.Title != "Buy milk" {
        t.Errorf("Expected title 'Buy milk', got '%s'", todo.Title)
    }
}

func TestCreateTodoUsecase_Execute_EmptyTitle(t *testing.T) {
    mock := &MockRepository{}
    usecase := NewCreateTodoUsecase(mock)

    _, err := usecase.Execute(context.Background(), "")

    if err == nil {
        t.Error("Expected error for empty title")
    }
}
```

#### 3.1.2 Verify RED: テストが失敗することを確認

```bash
go test ./internal/usecase -v
# コンパイルエラー: NewCreateTodoUsecase が未定義
```

#### 3.1.3 GREEN: テストを通す最小限のコードを書く

**TODO:**
- [ ] `internal/usecase/create_todo.go` を作成
- [ ] `CreateTodoUsecase` 構造体定義
- [ ] `Execute` メソッド実装

**実装:**
```go
package usecase

import (
    "context"
    "time"

    "github.com/k98a73/go-todo/internal/domain"
)

type CreateTodoUsecase struct {
    repo domain.IRepository
}

func NewCreateTodoUsecase(repo domain.IRepository) *CreateTodoUsecase {
    return &CreateTodoUsecase{repo: repo}
}

func (u *CreateTodoUsecase) Execute(ctx context.Context, title string) (*domain.Todo, error) {
    now := time.Now()
    todo := &domain.Todo{
        Title:     title,
        Completed: false,
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err := domain.ValidateTodo(todo); err != nil {
        return nil, err
    }

    if err := u.repo.Create(ctx, todo); err != nil {
        return nil, err
    }

    return todo, nil
}
```

#### 3.1.4 Verify GREEN

```bash
go test ./internal/usecase -v
# PASS
```

**学習ポイント:**
- 依存注入パターン
- モックを使った単体テスト
- ユースケースの責務

### 3.2 List ユースケース（TDD）

**同じサイクルを繰り返す:**
1. RED: `list_todo_test.go` でテストを書く
2. Verify RED: 失敗を確認
3. GREEN: `list_todo.go` で実装
4. Verify GREEN: 成功を確認

### 3.3 FindByID ユースケース（TDD）

**同じサイクルを繰り返す**

### 3.4 Update ユースケース（TDD）

**同じサイクルを繰り返す**

### 3.5 Delete ユースケース（TDD）

**同じサイクルを繰り返す**

---

## 段階4: インフラストラクチャ層実装（TDD）

### 4.1 JSON ファイルストレージの実装

#### 4.1.1 RED: 失敗するテストを書く

**TODO:**
- [ ] `internal/infra/storage/file_storage_test.go` を作成
- [ ] 一時ファイルを使用したテストを書く

**テストを先に書く:**
```go
package storage

import (
    "context"
    "os"
    "testing"

    "github.com/k98a73/go-todo/internal/domain"
)

func TestFileRepository_Create(t *testing.T) {
    tmpfile, err := os.CreateTemp("", "todo*.json")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())
    tmpfile.Write([]byte("[]"))
    tmpfile.Close()

    repo := NewFileRepository(tmpfile.Name())
    todo := &domain.Todo{Title: "Buy milk"}

    err = repo.Create(context.Background(), todo)

    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if todo.ID == 0 {
        t.Error("Expected ID to be assigned")
    }
}

func TestFileRepository_List(t *testing.T) {
    tmpfile, err := os.CreateTemp("", "todo*.json")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())
    tmpfile.Write([]byte(`[{"id":1,"title":"Test","completed":false}]`))
    tmpfile.Close()

    repo := NewFileRepository(tmpfile.Name())

    todos, err := repo.List(context.Background())

    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if len(todos) != 1 {
        t.Errorf("Expected 1 todo, got %d", len(todos))
    }
}
```

#### 4.1.2 Verify RED

```bash
go test ./internal/infra/storage -v
# コンパイルエラー: NewFileRepository が未定義
```

#### 4.1.3 GREEN: 実装

**TODO:**
- [ ] `internal/infra/storage/file_storage.go` を作成
- [ ] `FileRepository` 構造体で `IRepository` を実装

#### 4.1.4 Verify GREEN

```bash
go test ./internal/infra/storage -v
# PASS
```

**学習ポイント:**
- `os.ReadFile` / `os.WriteFile`
- `json.Marshal` / `json.Unmarshal`
- 一時ファイルを使ったテスト
- エラーハンドリング

---

## 段階5: HTTP 層実装（TDD）

### 5.1 HTTP ハンドラー実装

#### 5.1.1 RED: 失敗するテストを書く

**TODO:**
- [ ] `internal/infra/http/handler_test.go` を作成
- [ ] `httptest` を使用したテストを書く

**テストを先に書く:**
```go
package http

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestCreateTodoHandler(t *testing.T) {
    handler := NewTodoHandler(mockUsecase)

    body := strings.NewReader(`{"title": "Buy milk"}`)
    req, _ := http.NewRequest("POST", "/todo", body)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    handler.CreateTodo(w, req)

    if w.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", w.Code)
    }
}

func TestCreateTodoHandler_EmptyTitle(t *testing.T) {
    handler := NewTodoHandler(mockUsecase)

    body := strings.NewReader(`{"title": ""}`)
    req, _ := http.NewRequest("POST", "/todo", body)
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    handler.CreateTodo(w, req)

    if w.Code != http.StatusBadRequest {
        t.Errorf("Expected status 400, got %d", w.Code)
    }
}
```

#### 5.1.2 Verify RED → GREEN → Verify GREEN

**同じサイクルを繰り返す**

### 5.2 ルーティング設定

**TODO:**
- [ ] `cmd/main.go` でルーティング定義
- [ ] `http.ServeMux` 使用

```go
func main() {
    mux := http.NewServeMux()

    handler := newTodoHandler()
    mux.HandleFunc("POST /todo", handler.CreateTodo)
    mux.HandleFunc("GET /todo/list", handler.ListTodo)
    // ...

    http.ListenAndServe(":8080", mux)
}
```

---

## 段階6: 統合テスト・最適化

### 6.1-6.3 テスト実行・カバレッジ・手動テスト

詳細なテストコマンドと API 動作確認は、[COMMANDS.md](./COMMANDS.md#テスト実行) を参照してください。

- テスト実行: `go test ./...`
- カバレッジ測定: `go test -cover ./...`
- API 動作確認: curl コマンドで実行

### 6.4 エラーハンドリング改善（TDD）

**TODO:**
- [ ] エラーレスポンス形式のテストを書く
- [ ] テストに基づいて実装

**参考資料:** [ERROR_HANDLING.md](./ERROR_HANDLING.md)

### 6.5 パフォーマンス・セキュリティ対応

**TODO:**
- [ ] リクエストサイズ制限のテストを書く
- [ ] CORS 設定のテストを書く
- [ ] ログ出力の整備

---

## TDD チェックリスト

各機能実装時に確認：

- [ ] 失敗するテストを先に書いた
- [ ] テストが正しい理由で失敗することを確認した
- [ ] テストを通す最小限のコードを書いた
- [ ] すべてのテストが通ることを確認した
- [ ] コードをリファクタリングした（テストはグリーンを維持）
- [ ] エッジケースとエラーケースをカバーした

**すべてにチェックが入らない場合、TDDをスキップしています。やり直してください。**

---

## 学習マイルストーン

| 段階 | 達成内容              | 学習成果                           |
| ---- | --------------------- | ---------------------------------- |
| 1    | 環境構築、Go 基礎学習 | 言語の基本概念理解                 |
| 2    | domain 層完成（TDD）  | struct, interface, テスト駆動開発  |
| 3    | usecase 層完成（TDD） | ビジネスロジック、DI、モックテスト |
| 4    | infra 層完成（TDD）   | I/O 操作、JSON 処理、統合テスト    |
| 5    | HTTP 層完成（TDD）    | Web API 実装、httptest             |
| 6    | 全テスト＆最適化      | カバレッジ測定、エラーハンドリング |

---

## 参考資料

- [COMMANDS.md](./COMMANDS.md) - コマンドリファレンス
- [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) - ファイル構成
- [GO_BASICS.md](./GO_BASICS.md) - 言語の基礎
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 設計思想
- [TESTING.md](./TESTING.md) - テスト戦略
- [API_SPEC.md](./API_SPEC.md) - API 仕様
- [TDD スキル](../../.agent/skills/test-driven-development/SKILL.md) - TDD ガイド

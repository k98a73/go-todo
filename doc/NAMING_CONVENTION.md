# ネーミング規則

## 基本原則

- **単数形 / `list` 接尾辞**
  - 全ての名称は単数形、または単数形に `list` を付ける
  - 例：`todo`、`todoList` （`todos` は使わない）
  - 理由：英語の複数形が不規則なため、`s` の有無では見落としやすいため

## ディレクトリとファイル

- **小文字**
- **ケバブケース**
- 例：`service`、`handler`、`user-service.go`、`todo-handler.go`

## Go言語 - パッケージ名

- **小文字**
- **単語は繋げない**（ケバブケースやスネークケースは使わない）
- **単数形**
- 例：`todo`、`user`、`service`
  - 避けるべき：`todos`、`todo-service`、`todo_service`

## Go言語 - 変数名・関数名

### 変数名

- **キャメルケース**
- **単数形または `list` 接尾辞**
- **Exported（大文字始まり）と非Exported（小文字始まり）で使い分け**

例（非Exported）：
```go
var todoList []Todo
var currentUser User
var isActive bool
var maxRetryCount int
```

例（Exported）：
```go
var TodoList []Todo
var CurrentUser User
var IsActive bool
var MaxRetryCount int
```

### 関数名・メソッド名

- **キャメルケース**
- **動詞で始まる**（可能なら）
- **Exported（大文字始まり）と非Exported（小文字始まり）で使い分け**

例（非Exported）：
```go
func createTodo(ctx context.Context, title string) (*Todo, error)
func getTodoList(ctx context.Context, userID string) ([]*Todo, error)
func validateTodoInput(input *TodoInput) error
```

例（Exported）：
```go
func CreateTodo(ctx context.Context, title string) (*Todo, error)
func GetTodoList(ctx context.Context, userID string) ([]*Todo, error)
func ValidateTodoInput(input *TodoInput) error
```

## Go言語 - 構造体名・インターフェース名

- **キャメルケース（パスカルケース）**
- **大文字始まり**（Exported）
- **単数形**
- **接尾辞**：
  - 構造体：なし、または意図が明確な場合のみ（例：`Config`、`Repository`）
  - インターフェース：`er` で終わる（例：`Reader`、`Writer`）

例：
```go
type Todo struct {
    ID    string
    Title string
}

type TodoRepository interface {
    GetByID(ctx context.Context, id string) (*Todo, error)
    List(ctx context.Context) ([]*Todo, error)
}

type TodoService struct {
    repo TodoRepository
}

type TodoHandler struct {
    service TodoService
}
```

## Go言語 - 定数名

- **キャメルケース（パスカルケース）**
- **大文字始まり**（Exported）
- **定数グループは `const (...)` で囲む**

例：
```go
const (
    DefaultPageSize = 20
    MaxPageSize     = 100
)

const DefaultTimeout = 30 * time.Second
```

## Go言語 - エラーハンドリング

- 変数名は `err` で統一
- エラー定義は `Err` で始まる

例：
```go
var (
    ErrTodoNotFound = errors.New("todo not found")
    ErrInvalidInput = errors.New("invalid input")
)

if err != nil {
    return fmt.Errorf("failed to fetch todo: %w", err)
}
```

## 引数・戻り値の命名

### Context

- 最初の引数は常に `ctx context.Context`

```go
func CreateTodo(ctx context.Context, title string) (*Todo, error)
```

### スライス・配列

- `list` 接尾辞、または `Slice` は使わない
- 単数形で型名を示す、または明確な名前をつける

例：
```go
func ProcessTodos(todos []*Todo) error
func FilterByStatus(items []*Todo, status string) []*Todo
```

### Map・辞書

- 適切な変数名で、キーと値の関係を明確にする

例：
```go
userByID := make(map[string]*User)
statusCounts := make(map[string]int)
```

## テストコード

- テストファイル：`*_test.go`
- テスト関数：`Test` で始まる、テスト対象を含める

例：
```go
func TestCreateTodo(t *testing.T)
func TestGetTodoList_EmptyResult(t *testing.T)
func TestValidateTodoInput_InvalidTitle(t *testing.T)
```

## JSON/データモデルのタグ

- **snake_case** で定義（JSON API標準）

例：
```go
type Todo struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    CreatedAt time.Time `json:"created_at"`
}
```

## URL・APIエンドポイント

- **ケバブケース**
- **リソース名は単数形（RESTful慣例）**

例：
```
/api/v1/todo/{id}
/api/v1/todo
/api/v1/user/{id}/todo
```

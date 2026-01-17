# アーキテクチャ設計

## 採用する構造：**シンプルなレイヤード・アーキテクチャ**

Go初心者向けに、3層構造を採用します。各層が責務を持ち、依存関係が一方向です。

```
┌─────────────────────────────────┐
│   HTTP Handler Layer            │ ← クライアントからのリクエスト受け取り
│   (handlers/handler.go)         │ ← JSON解析、レスポンス返却
└──────────────┬──────────────────┘
               ↓
┌─────────────────────────────────┐
│   Business Logic Layer          │ ← CRUD処理のロジック
│   (services/todo-service.go)    │ ← バリデーション
└──────────────┬──────────────────┘
               ↓
┌─────────────────────────────────┐
│   Data Access Layer             │ ← ファイル読み書き
│   (repository/todo-repository.go)│ ← データの永続化・取得
└─────────────────────────────────┘
```

## 各層の責務

### 1. HTTP Handler Layer（ハンドラー層）
- **ファイル**: `handlers/handler.go`
- **役割**:
  - HTTPリクエストを受け取る
  - JSON → Go struct に解析
  - Go struct → JSON に変換
  - HTTPステータスコードを返却
  - **Service層のメソッドを呼び出す**
- **依存関係**: Service層に依存

**例**:
```go
func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    var todo models.Todo
    json.NewDecoder(r.Body).Decode(&todo)
    
    createdTodo, err := h.service.CreateTodo(todo)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    json.NewEncoder(w).Encode(createdTodo)
}
```

### 2. Business Logic Layer（ビジネスロジック層）
- **ファイル**: `services/todo-service.go`
- **役割**:
  - CRUD処理の実装
  - バリデーション（titleが空でないか等）
  - タイムスタンプの更新
  - **Repository層のメソッドを呼び出す**
- **依存関係**: Repository層に依存

**例**:
```go
func (s *TodoService) CreateTodo(todo models.Todo) (models.Todo, error) {
    if todo.Title == "" {
        return models.Todo{}, errors.New("title is required")
    }
    
    todo.CreatedAt = time.Now()
    todo.UpdatedAt = time.Now()
    
    return s.repo.Save(todo)
}
```

### 3. Data Access Layer（データアクセス層）
- **ファイル**: `repository/todo-repository.go`
- **役割**:
  - JSONファイルの読み書き
  - データの永続化（保存）
  - データの取得（読み込み）
  - **他の層に依存しない**
- **依存関係**: なし

**例**:
```go
func (r *TodoRepository) Save(todo models.Todo) (models.Todo, error) {
    todos, _ := r.Load()
    todos = append(todos, todo)
    
    data, _ := json.Marshal(todos)
    os.WriteFile(r.filePath, data, 0644)
    
    return todo, nil
}
```

## メリット

- **関心の分離**: 各層がやることが明確
- **テストしやすい**: 各層を独立してテスト可能
- **拡張性**: ファイル保存 → DB移行も容易
- **初心者向け**: 構造が理解しやすい

## 設計パターン

この構造は「Clean Architecture」「Hexagonal Architecture」の簡易版です。
より複雑なアプリケーションでは、さらに層を分けたり、インターフェースを活用したりします。

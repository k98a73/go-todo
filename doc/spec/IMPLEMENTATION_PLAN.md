# 実装計画

Go Todoアプリケーション開発の実装順序と各段階での学習ポイントをまとめます。

## 全体フロー

```
段階1: 基礎準備
  ↓
段階2: ドメイン層実装
  ↓
段階3: ユースケース層実装
  ↓
段階4: インフラストラクチャ層実装
  ↓
段階5: HTTP層実装
  ↓
段階6: テスト・最適化
```

---

## 段階1: 基礎準備（学習）

### 1.1 開発環境セットアップ

**実施内容:**
- Go のインストール
- リポジトリの作成
- リポジトリのクローン
- プロジェクトディレクトリ初期化
- `go.mod` ファイル生成
- mainブランチの作成

**学習ポイント:**
- Go のパッケージシステムの理解
- Module の基本
- Jujutsu（jj）によるバージョン管理

### 1.2 Go 言語の基礎学習

**学習内容:**
- struct（構造体）
- interface（インターフェース）
- time.Time
- context.Context
- error インターフェース
- defer
- ポインタ

**参考資料:** [GO_BASICS.md](./GO_BASICS.md)

**実装コマンド:**
詳細なコマンドリファレンスについては、[COMMANDS.md](./COMMANDS.md#初期セットアップコマンド) を参照してください。

---

## 段階2: ドメイン層実装

### 2.1 TODO 構造体の定義

**実施内容:**
- `internal/domain/entity.go` で `Todo` 構造体を定義
- JSON マーシャリング対応

**期待される実装:**
```go
type Todo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**学習ポイント:**
- struct フィールド定義
- struct タグ（JSON マーシャリング）
- time.Time の使用

### 2.2 バリデーション関数の実装

**実施内容:**
- Title のバリデーション（空でない、最大長チェック）
- 専用関数を `internal/domain/entity.go` に定義

**期待される実装:**
```go
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

**学習ポイント:**
- エラー返却パターン
- バリデーション設計

### 2.3 リポジトリインターフェース定義

**実施内容:**
- `IRepository` インターフェース定義
- 5つのメソッド（Create, List, FindByID, Update, Delete）

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

### 2.4 ドメイン層テスト

**実施内容:**
- `internal/domain/entity_test.go` を作成
- バリデーション関数のテスト実装

**テストケース例:**
- 正常な Title でテスト
- 空の Title でテスト
- 長すぎる Title でテスト

**学習ポイント:**
- ユニットテスト基本
- テーブル駆動テスト
- `testing.T` の使用方法

---

## 段階3: ユースケース層実装

### 3.1 Create ユースケース

**実施内容:**
- `internal/usecase/create_todo.go`
- `CreateTodoUsecase` 構造体定義
- リポジトリを DI（依存注入）で受け取る

**期待される実装:**
```go
type CreateTodoUsecase struct {
    repo domain.IRepository
}

func (u *CreateTodoUsecase) Execute(ctx context.Context, title string) (*domain.Todo, error) {
    // バリデーション
    // リポジトリに保存
    // 結果返却
}
```

**学習ポイント:**
- 依存注入パターン
- ユースケースの責務

### 3.2 List ユースケース

**実施内容:**
- `internal/usecase/list_todo.go`
- リポジトリから全 TODO を取得

**学習ポイント:**
- スライスの扱い
- インターフェース経由でのデータ取得

### 3.3 FindByID ユースケース

**実施内容:**
- `internal/usecase/find_todo.go`
- ID で単一 TODO を取得

### 3.4 Update ユースケース

**実施内容:**
- `internal/usecase/update_todo.go`
- 既存 TODO の更新
- バリデーション後に保存

### 3.5 Delete ユースケース

**実施内容:**
- `internal/usecase/delete_todo.go`
- ID で TODO を削除

### 3.6 ユースケース層テスト

**実施内容:**
- 各ユースケースのテスト
- **モックリポジトリ**を使用した単体テスト

**テスト例:**
```go
type MockRepository struct { /* ... */ }

func TestCreateTodoUsecase(t *testing.T) {
    mock := &MockRepository{}
    usecase := &CreateTodoUsecase{repo: mock}
    
    // テスト実行...
}
```

**学習ポイント:**
- モック実装
- インターフェースの活用

---

## 段階4: インフラストラクチャ層実装

### 4.1 JSON ファイルストレージの実装

**実施内容:**
- `internal/infra/storage/file_storage.go`
- `FileRepository` 構造体で `IRepository` を実装

**実装する機能:**
- JSON ファイル読み込み（`List`）
- JSON ファイル書き込み（`Create`, `Update`, `Delete`）
- ファイル I/O エラーハンドリング

**期待されるファイル構造:**
```json
[
  {
    "id": 1,
    "title": "Buy milk",
    "completed": false,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
]
```

**学習ポイント:**
- `os.ReadFile` / `os.WriteFile`
- `json.Marshal` / `json.Unmarshal`
- `defer` によるファイルクローズ
- エラーハンドリング

### 4.2 インフラ層テスト

**実施内容:**
- 一時ファイルを使用した統合テスト
- 実際のファイルI/O をテスト

```go
func TestFileRepositoryPersistence(t *testing.T) {
    tmpfile, _ := ioutil.TempFile("", "todo")
    defer os.Remove(tmpfile.Name())
    
    repo := NewFileRepository(tmpfile.Name())
    // テスト...
}
```

**学習ポイント:**
- 統合テスト設計
- 一時ファイル作成・削除

---

## 段階5: HTTP 層実装

### 5.1 HTTP ハンドラー実装

**実施内容:**
- `internal/infra/http/handler.go`
- `TodoHandler` 構造体定義
- 各メソッドでエンドポイント処理

**エンドポイント:**
- `POST /todo` → Create
- `GET /todo/list` → List
- `GET /todo/:id` → FindByID
- `PUT /todo/:id` → Update
- `DELETE /todo/:id` → Delete

**期待される実装:**
```go
type TodoHandler struct {
    usecase *usecase.CreateTodoUsecase
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    // JSON デコード
    // ユースケース呼び出し
    // JSON エンコード + レスポンス
}
```

**学習ポイント:**
- `http.ResponseWriter` / `*http.Request`
- JSON エンコード・デコード
- HTTP ステータスコード

### 5.2 ルーティング設定

**実施内容:**
- `cmd/main.go` でルーティング定義
- `http.HandleFunc` または `http.ServeMux` 使用

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

### 5.3 HTTP テスト

**実施内容:**
- `internal/infra/http/handler_test.go`
- `httptest` を使用したエンドポイントテスト

```go
func TestCreateTodoHandler(t *testing.T) {
    body := strings.NewReader(`{"title": "Buy milk"}`)
    req, _ := http.NewRequest("POST", "/todo", body)
    w := httptest.NewRecorder()
    
    handler.CreateTodo(w, req)
    
    if w.Code != http.StatusCreated {
        t.Errorf("Expected 201, got %d", w.Code)
    }
}
```

---

## 段階6: テスト・最適化

### 6.1-6.3 テスト実行・カバレッジ・手動テスト

詳細なテストコマンドと API 動作確認は、[COMMANDS.md](./COMMANDS.md#テスト実行) を参照してください。

- テスト実行: `go test ./...`
- カバレッジ測定: `go test -cover ./...`
- API 動作確認: curl コマンドで実行

### 6.4 エラーハンドリング改善

**実施内容:**
- エラーレスポンス形式の統一
- エラーコードの定義
- ロギング機能の追加

**参考資料:** [ERROR_HANDLING.md](./ERROR_HANDLING.md)

### 6.5 パフォーマンス・セキュリティ対応

- リクエストサイズ制限
- CORS 設定
- ログ出力の整備

---

## 学習マイルストーン

| 段階 | 達成内容 | 学習成果 |
|------|--------|--------|
| 1 | 環境構築、Go 基礎学習 | 言語の基本概念理解 |
| 2 | domain 層完成 | struct, interface, 設計思想 |
| 3 | usecase 層完成 | ビジネスロジック実装、DI パターン |
| 4 | infra 層完成 | I/O 操作、JSON 処理、統合テスト |
| 5 | HTTP 層完成 | Web API 実装、ルーティング |
| 6 | 全テスト＆最適化 | テスト設計、エラーハンドリング |

---

## 参考資料

- [COMMANDS.md](./COMMANDS.md) - コマンドリファレンス
- [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) - ファイル構成
- [GO_BASICS.md](./GO_BASICS.md) - 言語の基礎
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 設計思想
- [TESTING.md](./TESTING.md) - テスト戦略
- [API_SPEC.md](./API_SPEC.md) - API 仕様


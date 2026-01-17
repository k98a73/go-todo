# Go 学習ポイント

このプロジェクトで理解すべきGo言語の基礎知識をまとめます。

## struct（構造体）

### 定義方法

```go
type Todo struct {
    ID        int       // フィールド
    Title     string
    Completed bool
    CreatedAt time.Time
}
```

### メソッド

構造体にメソッドを定義できます：

```go
// レシーバーで構造体を指定
func (t *Todo) Complete() {
    t.Completed = true
}

func (t Todo) IsCompleted() bool {
    return t.Completed
}
```

**ポイント**
- `(t *Todo)` : ポインタレシーバー（構造体の値を変更可能）
- `(t Todo)` : 値レシーバー（構造体のコピーで処理、変更不可）

## interface（インターフェース）

### 定義方法

```go
type IRepository interface {
    Create(ctx context.Context, todo *Todo) error
    FindByID(ctx context.Context, id int) (*Todo, error)
    List(ctx context.Context) ([]*Todo, error)
    Update(ctx context.Context, todo *Todo) error
    Delete(ctx context.Context, id int) error
}
```

### ポイント

- インターフェースは**メソッドシグネチャ**を定義するだけ
- 実装側で明示的に実装を宣言する必要がない（暗黙的実装）
- ダックタイピングに似た仕組み

## time.Time

### 基本的な使い方

```go
now := time.Now()                          // 現在時刻
created := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

year := now.Year()                         // 2024
month := now.Month()                       // January
day := now.Day()                           // 15
hour := now.Hour()                         // 0

formatted := now.Format("2006-01-02 15:04:05")  // フォーマット
```

**重要**: Goの時刻フォーマットは `2006-01-02 15:04:05` という特定の値を使う

### JSON との連携

```go
type Todo struct {
    CreatedAt time.Time `json:"created_at"`
}

// JSON出力時は ISO 8601 形式（例: "2024-01-15T10:30:00Z"）
data, _ := json.Marshal(todo)
```

## context.Context

### 用途

- キャンセル信号の伝播
- タイムアウト管理
- リクエスト固有の値の保持

### 基本的な使い方

```go
// タイムアウト付きコンテキスト
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := repository.FindByID(ctx, 1)
```

### HTTP ハンドラーでの使い方

```go
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // リクエストからコンテキスト取得
    
    todo, err := h.usecase.FindByID(ctx, id)
    // ...
}
```

## error インターフェース

### エラー処理の基本

```go
type error interface {
    Error() string
}
```

### カスタムエラー

```go
type ValidationError struct {
    Message string
    Field   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s at field %s", e.Message, e.Field)
}
```

### エラーチェック

```go
if err != nil {
    log.Printf("error occurred: %v", err)
    return err
}
```

## defer（遅延実行）

### 用途

- リソースのクリーンアップ（ファイルクローズ等）
- パニック時の処理

### 基本的な使い方

```go
file, err := os.Open("data.json")
if err != nil {
    return err
}
defer file.Close()  // 関数を抜ける際に実行

// ファイル操作
```

## ゴルーチンとチャネル（中級）

プロジェクト初期段階では不要ですが、今後の拡張で必要になる可能性があります。

```go
// ゴルーチン：並行処理
go doSomething()

// チャネル：ゴルーチン間の通信
results := make(chan string)
go func() {
    results <- "done"
}()
message := <-results
```

## ポインタ

### 基本概念

```go
var x int = 42
var p *int = &x      // アドレス取得
value := *p           // デリファレンス（値を取得）
```

### 関数の引数

```go
func Modify(t *Todo) {      // ポインタ受け取り → 元の値が変更される
    t.Completed = true
}

func Read(t Todo) {         // 値受け取り → コピーなので元の値は変更されない
    _ = t.Completed
}
```

## まとめ

このプロジェクトで特に重要な概念：

1. **struct** → TODO のデータ構造
2. **interface** → リポジトリパターン
3. **time.Time** → 作成日時の管理
4. **context.Context** → リクエスト処理
5. **error** → エラーハンドリング
6. **defer** → ファイルクローズ
7. **ポインタ** → 値の変更・パフォーマンス

詳細は[Go公式ドキュメント](https://golang.org/doc/)を参照してください。

# ファイル保存仕様

## JSON形式での保存

### ファイルパス
```
./data/todo.json
```

### ファイル形式

**複数TODOの場合**:
```json
[
  {
    "id": 1,
    "title": "Go学習",
    "description": "Clean Architectureを学ぶ",
    "due_date": "2026-02-28T23:59:59Z",
    "completed": false,
    "created_at": "2026-01-17T10:00:00Z",
    "updated_at": "2026-01-17T10:00:00Z"
  },
  {
    "id": 2,
    "title": "テスト書き",
    "description": "",
    "due_date": "2026-02-01T00:00:00Z",
    "completed": true,
    "created_at": "2026-01-16T14:30:00Z",
    "updated_at": "2026-01-17T09:00:00Z"
  }
]
```

**TODOが0件の場合**:
```json
[]
```

### 文字コード
- **UTF-8**（Goのデフォルト）

### ファイルパーミッション
```go
os.WriteFile(filePath, data, 0644)
// 0644: オーナーは読み書き可、他は読み取り専用
```

---

## ファイルI/O処理

### 実装で使う標準パッケージ

```go
import (
    "encoding/json"  // JSONのマーシャル/アンマーシャル
    "os"             // ファイル操作
)
```

---

## ファイル読み込み

### 基本的な流れ

```go
// 1. ファイルを読み込む
data, err := os.ReadFile("./data/todo.json")
if err != nil {
    // ファイルが存在しない場合等
    return
}

// 2. JSON を Go struct にデコード
var todoList []models.Todo
err = json.Unmarshal(data, &todoList)
if err != nil {
    // JSON形式が不正
    return
}

// 3. 使用
for _, todo := range todoList {
    fmt.Println(todo.Title)
}
```

### ファイルが存在しない場合の処理

```go
data, err := os.ReadFile("./data/todo.json")
if err != nil {
    if os.IsNotExist(err) {
        // ファイルが存在しない場合
        return []models.Todo{}, nil  // 空のスライスを返す
    }
    // その他のエラー
    return nil, err
}
```

---

## ファイル書き込み

### 基本的な流れ

```go
// 1. Go struct を JSON にエンコード
data, err := json.Marshal(todoList)
if err != nil {
    return err
}

// 2. ファイルに書き込む
err = os.WriteFile("./data/todo.json", data, 0644)
if err != nil {
    return err
}

return nil
```

### ディレクトリの自動作成

```go
// data ディレクトリがない場合は作成
err := os.MkdirAll("./data", 0755)
if err != nil {
    return err
}

// ファイル書き込み
os.WriteFile("./data/todo.json", data, 0644)
```

---

## JSON マーシャル/アンマーシャル

### タグを使ったカスタマイズ

```go
type Todo struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    DueDate     time.Time `json:"due_date"`  // ケーバブケースに変換
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### JSONへの自動変換

```go
// time.Time は RFC3339 形式に自動変換される
todo := models.Todo{
    DueDate: time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC),
}

data, _ := json.Marshal(todo)
// Output: {"due_date":"2026-02-28T23:59:59Z",...}
```

### JSON からの自動変換

```go
// RFC3339 形式の文字列は自動的に time.Time に変換される
json_str := `{"due_date":"2026-02-28T23:59:59Z"}`

var todo models.Todo
json.Unmarshal([]byte(json_str), &todo)
// todo.DueDate は time.Time 型に変換済み
```

---

## パフォーマンス上の考慮

### 現在の設計の制限
- **ファイル全体を読み込み**: 毎回全データを読む
- **ファイル全体を上書き**: 変更時は全データを書き直す
- **スケーラビリティ**: 数千件程度まで問題なし（数万件以上は遅い）

### 将来的な改善案
- **インメモリキャッシュ**: ファイル読み込み回数を削減
- **データベース使用**: SQLite や PostgreSQL への移行
- **インデックス**: ID検索を高速化

---

## ファイルロック処理（将来の実装）

```go
// 複数のプロセスが同時にアクセスする場合はファイルロックが必要
// 現在の設計では考慮しない（単一プロセス想定）
```

---

## ファイル構造の初期化

### 初回起動時の処理

```go
// data/todo.json がない場合
if _, err := os.Stat("./data/todo.json"); os.IsNotExist(err) {
    // ファイルを作成
    os.MkdirAll("./data", 0755)
    os.WriteFile("./data/todo.json", []byte("[]"), 0644)
}
```

---

## テストでのファイル操作

```go
// テスト中は一時ファイルを使用
func TestSaveTodo(t *testing.T) {
    tmpFile := "/tmp/todo_test.json"
    repo := repository.NewTodoRepository(tmpFile)
    
    // テスト実行
    repo.Save(todo)
    
    // クリーンアップ
    os.Remove(tmpFile)
}
```

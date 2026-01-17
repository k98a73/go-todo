# データモデル定義

## TODO構造体（Go struct）

```go
type Todo struct {
    ID          int       // 一意識別子（自動採番）
    Title       string    // タイトル（必須）
    Description string    // 説明（オプション、空文字列可）
    DueDate     time.Time // 期日（日付型）
    Completed   bool      // 完了フラグ（デフォルト: false）
    CreatedAt   time.Time // 作成日時（自動生成）
    UpdatedAt   time.Time // 更新日時（自動生成）
}
```

## フィールド詳細

| フィールド | 型 | 説明 | 例 | 必須 |
|-----------|-----|------|----|----|
| ID | `int` | 一意識別子 | `1`, `2`, `3` | ✓ |
| Title | `string` | TODO のタイトル | `"Go学習"` | ✓ |
| Description | `string` | TODO の詳細説明 | `"Clean Architecture を学ぶ"` | ✗ |
| DueDate | `time.Time` | 期限日時 | `2026-02-28T23:59:59Z` | ✓ |
| Completed | `bool` | 完了状態 | `false`, `true` | ✗ |
| CreatedAt | `time.Time` | 作成日時 | `2026-01-17T10:00:00Z` | ✓ |
| UpdatedAt | `time.Time` | 最終更新日時 | `2026-01-17T15:30:00Z` | ✓ |

## Go の型について

### int
- **説明**: 整数型
- **使用例**: `var id int = 1`
- **値域**: システムに依存（64bit なら -9223372036854775808 〜 9223372036854775807）
- **Goで推奨**: 単に `int` とだけ書く（整数が必要なら自動的に調整）

### string
- **説明**: 文字列型
- **使用例**: `var title string = "Go学習"`
- **特徴**: 変更不可（イミュータブル）
- **初期値**: `""` (空文字列)

### bool
- **説明**: 真偽値型
- **使用例**: `var completed bool = false`
- **値**: `true` または `false`
- **初期値**: `false`

### time.Time
- **説明**: Goの標準日付型
- **使用例**:
  ```go
  import "time"
  
  now := time.Now()                          // 現在の日時
  due := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)  // 指定日時を作成
  parsed, _ := time.Parse(time.RFC3339, "2026-02-28T23:59:59Z")  // 文字列をパース
  ```
- **JSON化**: RFC3339形式 (`2026-02-28T23:59:59Z`) に自動変換

### slice（スライス）
- **説明**: 可変長配列
- **使用例**: `[]Todo` （TODOのスライス）
  ```go
  var todos []Todo
  todos = append(todos, newTodo)  // 追加
  ```

## JSON表現

### 構造体 → JSON への変換例

```go
todo := models.Todo{
    ID:          1,
    Title:       "Go学習",
    Description: "Clean Architectureを学ぶ",
    DueDate:     time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC),
    Completed:   false,
    CreatedAt:   time.Now(),
    UpdatedAt:   time.Now(),
}

data, _ := json.Marshal(todo)
// Output:
// {
//   "ID": 1,
//   "Title": "Go学習",
//   "Description": "Clean Architectureを学ぶ",
//   "DueDate": "2026-02-28T23:59:59Z",
//   "Completed": false,
//   "CreatedAt": "2026-01-17T10:00:00Z",
//   "UpdatedAt": "2026-01-17T10:00:00Z"
// }
```

### JSONタグ（ケーバブケースでシリアライズする場合）

```go
type Todo struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    DueDate     time.Time `json:"due_date"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## 初期化方法

### リテラル記法
```go
todo := models.Todo{
    ID:          1,
    Title:       "Go学習",
    DueDate:     time.Now(),
    Completed:   false,
    CreatedAt:   time.Now(),
    UpdatedAt:   time.Now(),
}
```

### ゼロ値（初期化なし）
```go
var todo models.Todo
// ID: 0, Title: "", Description: "", DueDate: 0001-01-01T00:00:00Z, Completed: false, CreatedAt: 0001-01-01T00:00:00Z, UpdatedAt: 0001-01-01T00:00:00Z
```

## バリデーションルール

実装時に以下のルールで検証します：

1. **Title**: 必須、空文字列不可
2. **DueDate**: 必須、未来の日付推奨（過去でもOK）
3. **Description**: オプション（空文字列OK）
4. **Completed**: デフォルト `false`
5. **CreatedAt/UpdatedAt**: サーバー側で自動生成、上書き不可

# エラーハンドリング設計

## HTTPステータスコード

| コード | 名前 | 用途 | 例 |
|--------|------|------|-----|
| `200` | OK | リクエスト成功（GET, PUT, DELETE） | |
| `201` | Created | リソース作成成功（POST） | |
| `400` | Bad Request | クライアント側の入力エラー | titleが空文字列 |
| `404` | Not Found | リソースが見つからない | 存在しないIDにアクセス |
| `500` | Internal Server Error | サーバー内部エラー | ファイル読み書き失敗 |

---

## エラーレスポンスのフォーマット

すべてのエラーレスポンスは統一フォーマットで返します：

```json
{
  "error": "error_code",
  "message": "human-readable error message"
}
```

---

## 各エンドポイントのエラーパターン

### GET /todo
**成功時**: `200 OK`
```json
[...]
```

**失敗パターン**: なし（常に成功）

---

### GET /todo/:id
**成功時**: `200 OK`
```json
{...}
```

**失敗パターン1: IDが見つからない**
```
HTTPステータス: 404 Not Found
{
  "error": "not_found",
  "message": "todo with id 999 not found"
}
```

**失敗パターン2: IDが数値でない**
```
HTTPステータス: 400 Bad Request
{
  "error": "invalid_id",
  "message": "id must be a number"
}
```

---

### POST /todo
**成功時**: `201 Created`
```json
{...}
```

**失敗パターン1: titleが空文字列**
```
HTTPステータス: 400 Bad Request
{
  "error": "invalid_request",
  "message": "title is required"
}
```

**失敗パターン2: due_dateの形式が不正**
```
HTTPステータス: 400 Bad Request
{
  "error": "invalid_date",
  "message": "due_date must be in RFC3339 format (e.g., 2026-02-28T23:59:59Z)"
}
```

**失敗パターン3: JSONのデコードに失敗**
```
HTTPステータス: 400 Bad Request
{
  "error": "invalid_json",
  "message": "request body is not valid JSON"
}
```

**失敗パターン4: ファイル保存に失敗**
```
HTTPステータス: 500 Internal Server Error
{
  "error": "storage_error",
  "message": "failed to save todo"
}
```

---

### PUT /todo/:id
**成功時**: `200 OK`
```json
{...}
```

**失敗パターン1: IDが見つからない**
```
HTTPステータス: 404 Not Found
{
  "error": "not_found",
  "message": "todo with id 999 not found"
}
```

**失敗パターン2: リクエストボディが不正**
```
HTTPステータス: 400 Bad Request
{
  "error": "invalid_request",
  "message": "invalid request body"
}
```

**失敗パターン3: titleが空文字列に更新しようとした**
```
HTTPステータス: 400 Bad Request
{
  "error": "invalid_request",
  "message": "title cannot be empty"
}
```

---

### DELETE /todo/:id
**成功時**: `200 OK`
```json
{
  "message": "todo deleted successfully"
}
```

**失敗パターン1: IDが見つからない**
```
HTTPステータス: 404 Not Found
{
  "error": "not_found",
  "message": "todo with id 999 not found"
}
```

---

## Go での error 型と扱い方

### error 型の基本

```go
// error は interface で、Error() メソッドを実装すればOK
type error interface {
    Error() string
}

// 通常はerrors.New()で作成
import "errors"

err := errors.New("title is required")
if err != nil {
    // エラー処理
}
```

### 複数の戻り値とエラーチェック

```go
// 通常、GoではエラーはLast(最後)の戻り値
todo, err := service.GetTodo(1)
if err != nil {
    // エラー処理
    return
}
// 成功時の処理
```

### HTTPハンドラーでのエラー処理

```go
func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")  // URLパラメータを取得
    
    todo, err := h.service.GetTodo(id)
    if err != nil {
        // エラー時はHTTPレスポンスで返す
        response := map[string]string{
            "error": "not_found",
            "message": err.Error(),
        }
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(response)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(todo)
}
```

---

## バリデーションエラー vs 実行時エラー

### バリデーションエラー（400 Bad Request）
- クライアント側のミス
- リクエストボディが不正
- 例: titleが空、due_dateの形式が不正

### 実行時エラー（500 Internal Server Error）
- サーバー側のミス
- データベース・ファイルのI/Oエラー
- 予期しない例外
- 例: ファイル保存失敗、ファイル読み込み失敗

---

## ロギング設計

```go
// DEBUGレベル: 処理の詳細
log.Printf("DEBUG: getting todo with id: %d", id)

// INFOレベル: 処理の成功
log.Printf("INFO: todo created with id: %d", todo.ID)

// ERRORレベル: エラー発生（スタックトレースは不要）
log.Printf("ERROR: failed to save todo: %v", err)

// 注意: クライアントにはエラー詳細を返さない（セキュリティのため）
```

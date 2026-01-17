# API 仕様

## エンドポイント一覧

### 全TODOを取得
- **メソッド**: `GET`
- **パス**: `/todo/list`
- **説明**: 保存されている全TODOを取得

**リクエスト**:
```bash
curl http://localhost:8080/todo/list
```

**レスポンス（成功時）**:
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

**HTTPステータス**:
- `200 OK`: 成功

---

### 特定のTODOを取得
- **メソッド**: `GET`
- **パス**: `/todo/:id`
- **説明**: IDで指定したTODOを取得

**リクエスト**:
```bash
curl http://localhost:8080/todo/1
```

**レスポンス（成功時）**:
```json
{
  "id": 1,
  "title": "Go学習",
  "description": "Clean Architectureを学ぶ",
  "due_date": "2026-02-28T23:59:59Z",
  "completed": false,
  "created_at": "2026-01-17T10:00:00Z",
  "updated_at": "2026-01-17T10:00:00Z"
}
```

**レスポンス（失敗時: IDが見つからない）**:
```json
{
  "error": "not found",
  "message": "todo with id 999 not found"
}
```

**HTTPステータス**:
- `200 OK`: 成功
- `404 Not Found`: TODOが見つからない

---

### 新規TODO作成
- **メソッド**: `POST`
- **パス**: `/todo`
- **説明**: 新しいTODOを作成

**リクエスト**:
```bash
curl -X POST http://localhost:8080/todo \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go学習",
    "description": "Clean Architectureを学ぶ",
    "due_date": "2026-02-28T23:59:59Z"
  }'
```

**リクエストボディ（必須フィールド）**:
```json
{
  "title": "Go学習",
  "description": "Clean Architectureを学ぶ",
  "due_date": "2026-02-28T23:59:59Z"
}
```

**注意**: `ID`, `CreatedAt`, `UpdatedAt` はリクエストで指定不可（サーバー側で自動生成）

**レスポンス（成功時）**:
```json
{
  "id": 1,
  "title": "Go学習",
  "description": "Clean Architectureを学ぶ",
  "due_date": "2026-02-28T23:59:59Z",
  "completed": false,
  "created_at": "2026-01-17T10:00:00Z",
  "updated_at": "2026-01-17T10:00:00Z"
}
```

**レスポンス（失敗時: titleが空）**:
```json
{
  "error": "invalid request",
  "message": "title is required"
}
```

**HTTPステータス**:
- `201 Created`: 作成成功
- `400 Bad Request`: リクエストが不正

---

### TODOを更新
- **メソッド**: `PUT`
- **パス**: `/todo/:id`
- **説明**: IDで指定したTODOを更新

**リクエスト**:
```bash
curl -X PUT http://localhost:8080/todo/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go学習（更新）",
    "completed": true
  }'
```

**リクエストボディ（部分更新対応）**:
```json
{
  "title": "Go学習（更新）",
  "description": "Clean Architectureを学ぶ",
  "due_date": "2026-03-31T23:59:59Z",
  "completed": true
}
```

**注意**: 指定したフィールドのみ更新（`ID`, `CreatedAt` は上書き不可）

**レスポンス（成功時）**:
```json
{
  "id": 1,
  "title": "Go学習（更新）",
  "description": "Clean Architectureを学ぶ",
  "due_date": "2026-03-31T23:59:59Z",
  "completed": true,
  "created_at": "2026-01-17T10:00:00Z",
  "updated_at": "2026-01-17T15:30:00Z"
}
```

**HTTPステータス**:
- `200 OK`: 更新成功
- `400 Bad Request`: リクエストが不正
- `404 Not Found`: TODOが見つからない

---

### TODOを削除
- **メソッド**: `DELETE`
- **パス**: `/todo/:id`
- **説明**: IDで指定したTODOを削除

**リクエスト**:
```bash
curl -X DELETE http://localhost:8080/todo/1
```

**レスポンス（成功時）**:
```json
{
  "message": "todo deleted successfully"
}
```

**HTTPステータス**:
- `200 OK`: 削除成功
- `404 Not Found`: TODOが見つからない

---

## リクエスト/レスポンス仕様

### リクエストの日付形式
- **形式**: RFC3339 (ISO8601)
- **例**: `2026-02-28T23:59:59Z`
- **パース方法**:
  ```go
  t, _ := time.Parse(time.RFC3339, "2026-02-28T23:59:59Z")
  ```

### Content-Type
- **リクエスト**: `Content-Type: application/json`
- **レスポンス**: `Content-Type: application/json`

### 文字コード
- **UTF-8**

---

## エラーレスポンスのフォーマット

```json
{
  "error": "error_code",
  "message": "human-readable error message"
}
```

**例**:
```json
{
  "error": "invalid request",
  "message": "title is required"
}
```

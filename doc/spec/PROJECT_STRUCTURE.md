# プロジェクト構成

Go Todoアプリケーションのディレクトリ・ファイル構成を説明します。

## ディレクトリ構造

```
go-todo/
├── cmd/
│   └── main.go              # アプリケーションのエントリーポイント
├── internal/
│   ├── domain/              # ビジネスロジック層（3層アーキテクチャ）
│   │   ├── entity.go        # TODO構造体の定義
│   │   └── repository.go    # リポジトリインターフェース
│   ├── usecase/             # ユースケース層（ビジネスロジック）
│   │   ├── create_todo.go   # TODO作成ロジック
│   │   ├── list_todo.go    # TODO一覧取得ロジック
│   │   ├── update_todo.go   # TODO更新ロジック
│   │   └── delete_todo.go   # TODO削除ロジック
│   └── infra/               # インフラストラクチャ層（外部連携）
│       ├── http/            # HTTPサーバー・ハンドラー
│       │   ├── handler.go   # エンドポイントハンドラー
│       │   └── middleware.go # HTTPミドルウェア
│       └── storage/         # ストレージ層
│           └── file_storage.go # JSON ファイル保存実装
├── pkg/                     # 共通ユーティリティ
│   ├── logger/             # ロギング機能
│   ├── errors/             # エラーハンドリング
│   └── validator/          # バリデーション
├── doc/                    # ドキュメント
│   └── spec/               # 設計仕様書
├── go.mod                  # Go モジュールファイル
├── go.sum                  # モジュール依存関係ロック
└── README.md               # プロジェクト説明
```

## 各層の責務

### `cmd/main.go`
- アプリケーションのエントリーポイント
- サーバー起動、初期化処理

### `internal/domain/`
- **ビジネスロジック層（最も重要）**
- TODO構造体の定義
- リポジトリインターフェース
- ビジネスルール（バリデーション、制約）

### `internal/usecase/`
- **ユースケース層**
- ビジネスロジックの組み合わせ
- ドメイン層のリポジトリを使用
- HTTP層とドメイン層の橋渡し

### `internal/infra/`
- **インフラストラクチャ層**
- HTTP通信の実装
- ファイルI/O
- 外部システムとの連携

### `pkg/`
- 複数の層で共通利用するユーティリティ

## 命名規則

- **ファイル名**: `snake_case`（例: `create_todo.go`）
- **パッケージ名**: `lowercase`（例: `domain`, `usecase`）
- **構造体名**: `PascalCase`（例: `Todo`, `CreateTodoUsecase`）
- **関数名**: `PascalCase`（公開）、`camelCase`（非公開）
- **インターフェース名**: `I`をプリフィックスまたはサフィックスに付与（例: `IRepository`）

詳細は[NAMING_CONVENTION.md](../NAMING_CONVENTION.md)を参照してください。

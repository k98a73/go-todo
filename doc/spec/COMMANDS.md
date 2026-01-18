# コマンドリファレンス

## 初期セットアップコマンド

リポジトリの初期設定時に実行するコマンド一覧

```bash
# リポジトリ作成（初回のみ）
gh repo create k98a73/go-todo --public

# リポジトリクローン
ghq get https://github.com/k98a73/go-todo

# リポジトリに移動
cd $(ghq root)/github.com/k98a73/go-todo

# Go プロジェクト初期化
go mod init github.com/k98a73/go-todo

# mainブランチの作成
jj bookmark create main -r luotxono
```

---

## 開発中によく使うコマンド

### プロジェクト実行

```bash
# メインプログラム実行
go run cmd/main.go

# 特定ファイルで実行
go run cmd/main.go --flags
```

### テスト実行

```bash
# 全テスト実行
go test ./...

# 特定パッケージのテスト
go test ./internal/domain

# 詳細出力
go test -v ./...

# カバレッジ測定
go test -cover ./...

# HTMLカバレッジレポート生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### コード品質確認

```bash
# フォーマット（確認のみ）
go fmt ./...

# ベストプラクティスチェック
go vet ./...

# 依存関係を確認
go mod tidy

# 依存関係のグラフ表示
go mod graph
```

### ビルド

```bash
# バイナリビルド
go build -o bin/todo cmd/main.go

# クロスコンパイル
GOOS=darwin GOARCH=amd64 go build -o bin/todo-macos cmd/main.go
```

---

## デバッグ・検査コマンド

```bash
# パッケージ依存関係を表示
go list -m all

# ソースコード解析
go doc github.com/k98a73/go-todo/internal/domain

# 実行時メモリ調査
go tool pprof http://localhost:6060/debug/pprof/heap

# ベンチマークテスト
go test -bench=. ./...
```

---

## API 動作確認（curl）

```bash
# サーバー起動（別ターミナル）
go run cmd/main.go

# TODO作成
curl -X POST http://localhost:8080/todo \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy milk"}'

# TODO一覧取得
curl http://localhost:8080/todo/list

# 特定TODOを取得
curl http://localhost:8080/todo/1

# TODO更新
curl -X PUT http://localhost:8080/todo/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy milk and eggs", "completed": true}'

# TODO削除
curl -X DELETE http://localhost:8080/todo/1
```

---

## Tips

- **go mod tidy**: 不要な依存関係を削除し、必要なものを追加
- **go fmt**: 全てのGoファイルをフォーマット（VS Codeの保存時自動実行推奨）
- **go vet**: 一般的なバグを自動検出
- **-v フラグ**: テスト詳細出力時に有用

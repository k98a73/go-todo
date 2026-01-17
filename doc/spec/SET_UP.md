# セットアップ

## 基本要件

- Goはインストール済み
- データを保存したい
  - できるだけ簡易な方法で
  - ファイルでも可能
  - 現在のデータの状況が確認できれば良い
- フロントエンドは無し
- TODOアプリのAPIを作成
  - CRUD操作が可能
  - APIは叩いて確認（curlか何かしらのツールを使うか）
  - 登録するデータは、Go言語で扱う主要な型を網羅するようにする
    - 日付は絶対に入れる

---

## 各計画ドキュメント

詳細な設計・学習内容は以下のファイルで管理します：

| ファイル | 内容 |
|---------|------|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | 3層レイヤード・アーキテクチャの設計 |
| [DATA_MODEL.md](./DATA_MODEL.md) | TODO構造体とフィールド定義 |
| [API_SPEC.md](./API_SPEC.md) | エンドポイント仕様とリクエスト/レスポンス |
| [ERROR_HANDLING.md](./ERROR_HANDLING.md) | エラーハンドリング設計 |
| [STORAGE.md](./STORAGE.md) | JSON形式でのファイル保存方法 |
| [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) | ディレクトリ・ファイル構成 |
| [GO_BASICS.md](./GO_BASICS.md) | Go学習ポイント（struct、time.Time等） |
| [TESTING.md](./TESTING.md) | テスト戦略（ユニットテスト、curl） |
| [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md) | 実装の流れと順序 |

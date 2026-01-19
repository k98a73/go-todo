# 開発ワークフロー

## テスト駆動開発（TDD）採用

### 基本サイクル

1. **テスト作成（Red）**
   - 実装すべき動作をテストコードで表現
   - テストが失敗することを確認

2. **実装（Green）**
   - テストに合格する最もシンプルなコードを実装
   - 過度な設計・汎化は行わない

3. **整理（Refactor）**
   - テストが緑のまま、コードを整理
   - 重複を排除、可読性向上

詳細は [`doc/spec/TESTING.md`](doc/spec/TESTING.md) または `.agents/skills/test-driven-development/SKILL.md` を参照。

## 設計思想

### 3つのアプローチを段階的に採用

1. **TDD（Test-Driven Development）**
   - テスト優先で実装
   - 仕様を明確に

2. **DDD（Domain-Driven Design）**
   - ビジネスロジックを中心に設計
   - ドメインモデル構築

3. **Clean Architecture**
   - 層の分離（UI層、Application層、Domain層、Infrastructure層）
   - 依存関係の管理

詳細は [`doc/spec/ARCHITECTURE.md`](doc/spec/ARCHITECTURE.md) を参照。

## 実装計画

全体の段階的な実装計画は [`doc/spec/IMPLEMENTATION_PLAN.md`](doc/spec/IMPLEMENTATION_PLAN.md) を参照。

## コード規約

- [`doc/NAMING_CONVENTION.md`](doc/NAMING_CONVENTION.md) - 変数名、関数名、構造体名など
- [`doc/spec/ERROR_HANDLING.md`](doc/spec/ERROR_HANDLING.md) - エラーハンドリング規則
- [`doc/spec/GO_BASICS.md`](doc/spec/GO_BASICS.md) - Go言語の基礎概念

## バージョン管理

Jujutsu（jj）を使用。基本的なワークフロー：

```bash
# 変更を記録
jj commit -m "メッセージ"

# リモートにプッシュ
jj git push
```

詳細は `.agents/skills/jujutsu/SKILL.md` を参照。

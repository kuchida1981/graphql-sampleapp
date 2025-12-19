# Spec: Testing Policy Documentation

## ADDED Requirements

### Requirement: Testing Policy Document
プロジェクトルートに`TESTING.md`ドキュメントが存在し、テスト作成の方針とガイドラインを記載しなければならない (MUST exist and document)。

#### Scenario: Document exists in project root
- **Given**: プロジェクトルートディレクトリ
- **When**: テストポリシードキュメントを配置する
- **Then**: `/TESTING.md`ファイルが存在する

---

### Requirement: Testing Policy Content - Basic Principles
`TESTING.md`にはテスト作成の基本方針が明記されていなければならない (MUST document)。

#### Scenario: Document basic testing principles
- **Given**: `TESTING.md`ドキュメント
- **When**: 基本方針セクションを読む
- **Then**: 以下の方針が記載されている:
  - すべての公開メソッドにテストを書く
  - DBは常にモック化する
  - テーブル駆動テストを使用する
  - 正常系と異常系の両方をテストする

---

### Requirement: Testing Policy Content - Mocking Strategy
`TESTING.md`にはモッキング戦略が明記されていなければならない (MUST document)。

#### Scenario: Document database mocking
- **Given**: `TESTING.md`ドキュメント
- **When**: モッキング戦略セクションを読む
- **Then**: 以下が記載されている:
  - PostgreSQLモックには`go-sqlmock`を使用
  - Firestoreモックの作成方法
  - Repositoryインターフェースのモック方法

---

### Requirement: Testing Policy Content - Test Naming Convention
`TESTING.md`にはテスト命名規則が明記されていなければならない (MUST document)。

#### Scenario: Document test naming convention
- **Given**: `TESTING.md`ドキュメント
- **When**: 命名規則セクションを読む
- **Then**: `Test<FunctionName>_<Scenario>`形式の命名規則が記載されている
- **And**: 具体例が提供されている

---

### Requirement: Testing Policy Content - Table-Driven Tests
`TESTING.md`にはテーブル駆動テストの書き方が説明されていなければならない (MUST explain)。

#### Scenario: Document table-driven test pattern
- **Given**: `TESTING.md`ドキュメント
- **When**: テーブル駆動テストセクションを読む
- **Then**: `[]struct{name, input, want, wantErr}`パターンの説明がある
- **And**: 実際のコード例が提供されている

---

### Requirement: Testing Policy Content - Coverage Goals
`TESTING.md`にはカバレッジ方針が明記されていなければならない (MUST document)。

#### Scenario: Document coverage policy
- **Given**: `TESTING.md`ドキュメント
- **When**: カバレッジセクションを読む
- **Then**: 以下が記載されている:
  - 具体的なパーセンテージ目標は設定しない
  - 意味のあるカバレッジを優先する
  - `go test -cover ./...`で確認できることが説明されている

---

### Requirement: Testing Policy Content - Running Tests
`TESTING.md`にはテスト実行方法が記載されていなければならない (MUST document)。

#### Scenario: Document how to run tests
- **Given**: `TESTING.md`ドキュメント
- **When**: テスト実行セクションを読む
- **Then**: 以下のコマンドが記載されている:
  - `go test ./...` - すべてのテスト実行
  - `go test -v ./...` - 詳細出力付き実行
  - `go test -cover ./...` - カバレッジ付き実行
  - `go test ./internal/repository` - 特定パッケージのテスト

---

### Requirement: Testing Policy Content - CI/CD Integration
`TESTING.md`にはCI/CDでのテスト実行方法が記載されていなければならない (MUST document)。

#### Scenario: Document CI/CD test execution
- **Given**: `TESTING.md`ドキュメント
- **When**: CI/CDセクションを読む
- **Then**: 以下が記載されている:
  - すべてのPRでテストが自動実行されること
  - テスト失敗時のマージブロック
  - GitHub Actions (または同等) の設定ファイル参照

---

### Requirement: Testing Policy Content - Examples
`TESTING.md`には実際のテストコード例が含まれていなければならない (MUST include)。

#### Scenario: Document repository test example
- **Given**: `TESTING.md`ドキュメント
- **When**: サンプルコードセクションを読む
- **Then**: PostgreSQL Repositoryのテスト例が記載されている

#### Scenario: Document resolver test example
- **Given**: `TESTING.md`ドキュメント
- **When**: サンプルコードセクションを読む
- **Then**: GraphQL Resolverのテスト例が記載されている

---

### Requirement: Testing Policy Content - Dependencies
`TESTING.md`には必要なテストライブラリが明記されていなければならない (MUST list)。

#### Scenario: Document required testing libraries
- **Given**: `TESTING.md`ドキュメント
- **When**: 依存関係セクションを読む
- **Then**: 以下のライブラリが記載されている:
  - `testing` (標準ライブラリ)
  - `github.com/DATA-DOG/go-sqlmock`
  - `github.com/golang/mock/gomock` (オプション)

---

### Requirement: Japanese Language Documentation
`TESTING.md`は日本語で記述されていなければならない (MUST be written in Japanese) (プロジェクトの他のドキュメントと一貫性を保つため)。

#### Scenario: Document language is Japanese
- **Given**: `TESTING.md`ドキュメント
- **When**: ドキュメントを開く
- **Then**: すべての説明が日本語で記述されている
- **And**: コード例のコメントも日本語である

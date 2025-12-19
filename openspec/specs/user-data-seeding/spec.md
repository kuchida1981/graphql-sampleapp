# user-data-seeding Specification

## Purpose
TBD - created by archiving change add-postgresql-user-integration. Update Purpose after archive.
## Requirements
### Requirement: Database Schema Initialization
システムは初期化SQLスクリプトでusersテーブルを作成することをMUSTとする。

#### Scenario: usersテーブルの定義
- **WHEN** `scripts/init-postgres.sql` を確認する
- **THEN** `CREATE TABLE IF NOT EXISTS users` ステートメントが存在する
- **AND** `id VARCHAR(255) PRIMARY KEY` カラムが定義されている
- **AND** `name VARCHAR(255) NOT NULL` カラムが定義されている
- **AND** `email VARCHAR(255) NOT NULL UNIQUE` カラムが定義されている
- **AND** `created_at TIMESTAMP NOT NULL DEFAULT NOW()` カラムが定義されている

#### Scenario: インデックスの定義
- **WHEN** `scripts/init-postgres.sql` を確認する
- **THEN** `CREATE INDEX idx_users_email ON users(email)` が定義されている
- **AND** `CREATE INDEX idx_users_created_at ON users(created_at)` が定義されている

#### Scenario: べき等性
- **WHEN** 初期化スクリプトを複数回実行する
- **THEN** `IF NOT EXISTS` により既存テーブルは変更されない
- **AND** エラーが発生しない

### Requirement: Sample Data Seeding
システムはPostgreSQLに初期サンプルデータを投入する手段を提供することをMUSTとする。

#### Scenario: シードスクリプトの存在
- **WHEN** プロジェクトディレクトリを確認する
- **THEN** `scripts/seed-postgres.go` ファイルが存在する
- **AND** スクリプトは `DATABASE_URL` 環境変数を読み込む

#### Scenario: サンプルデータの投入
- **WHEN** `go run scripts/seed-postgres.go` を実行する
- **THEN** 少なくとも3つのUserレコードが `users` テーブルに作成される
- **AND** 各レコードは `id`, `name`, `email`, `created_at` フィールドを持つ
- **AND** 成功メッセージがコンソールに表示される

#### Scenario: べき等性
- **WHEN** シードスクリプトを複数回実行する
- **THEN** `INSERT ... ON CONFLICT (id) DO UPDATE` または同等の処理で既存データを更新する
- **AND** エラーが発生しない

#### Scenario: サンプルデータの内容
- **WHEN** シードされたデータを確認する
- **THEN** ユーザー名は実在しない架空の名前である（例: "Alice Smith", "Bob Johnson"）
- **AND** メールアドレスは一意である
- **AND** 作成日時は現在時刻またはダミーの過去日時である


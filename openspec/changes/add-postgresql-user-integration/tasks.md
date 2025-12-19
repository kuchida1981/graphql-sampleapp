# Tasks: PostgreSQL User Integration

このタスクリストは、ユーザーに見える進捗を提供する順序で整理されています。各タスクは小さく検証可能で、依存関係と並列化可能な作業を明示しています。

## Phase 1: Foundation (基盤構築)

### 1. Docker Compose PostgreSQL設定
**目的**: PostgreSQLコンテナを起動可能にする
**依存**: なし
**検証**: `docker compose up postgres` でPostgreSQLが起動し、`psql`で接続確認できる

- `docker-compose.yml`に`postgres`サービスを追加
  - イメージ: `postgres:16-alpine`
  - 環境変数: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
  - ポート: `5432:5432`
  - ボリューム: `postgres_data`, 初期化スクリプト用マウント
- `volumes`セクションに`postgres_data`を追加

### 2. PostgreSQL初期化スクリプト作成
**目的**: usersテーブルを自動作成する
**依存**: タスク1（Docker Compose設定）
**検証**: PostgreSQLコンテナ起動後、`users`テーブルが存在する

- `scripts/init-postgres.sql`を作成
  - `CREATE TABLE IF NOT EXISTS users`
  - カラム: `id`, `name`, `email`, `created_at`
  - `email`にUNIQUE制約
  - インデックス: `idx_users_email`, `idx_users_created_at`

### 3. Userドメインモデル定義
**目的**: Userエンティティを定義する
**依存**: なし（並列化可能）
**検証**: `internal/domain/user.go`がコンパイル可能

- `internal/domain/user.go`を作成
  - `User`構造体定義（ID, Name, Email, CreatedAt）
  - データベースタグは不要（手動マッピング）

## Phase 2: Data Access Layer (データアクセス層)

### 4. PostgreSQLクライアント初期化
**目的**: PostgreSQL接続を確立する
**依存**: タスク1（PostgreSQLコンテナ）
**検証**: アプリケーション起動時にPostgreSQLに接続成功ログが出る

- `internal/postgres/client.go`を作成
  - `NewClient(ctx context.Context, connStr string) (*sql.DB, error)`
  - 環境変数`DATABASE_URL`から接続文字列を取得
  - `Ping()`で接続確認
  - エラーハンドリングとログ出力
- `go.mod`に依存追加: `github.com/jackc/pgx/v5/stdlib`

### 5. UserRepositoryインターフェース定義
**目的**: リポジトリパターンを定義する
**依存**: タスク3（Userドメインモデル）
**検証**: インターフェースがコンパイル可能

- `internal/repository/user.go`を作成
  - `UserRepository`インターフェース定義
  - メソッド: `List(ctx)`, `GetByID(ctx, id)`

### 6. PostgresUserRepository実装
**目的**: PostgreSQLからUserデータを取得する
**依存**: タスク4（PostgreSQLクライアント）、タスク5（インターフェース）
**検証**: ユニットテスト可能（または手動でシードデータ取得確認）

- `internal/repository/postgres_user.go`を作成
  - `PostgresUserRepository`構造体
  - コンストラクタ: `NewPostgresUserRepository(db *sql.DB)`
  - `List()`実装: `SELECT * FROM users ORDER BY created_at DESC`
  - `GetByID()`実装: `SELECT * FROM users WHERE id = $1`
  - エラーハンドリング: `sql.ErrNoRows`を適切に処理
  - ログ出力

## Phase 3: GraphQL Integration (GraphQL統合)

### 7. GraphQLスキーマ拡張
**目的**: User型とクエリを追加する
**依存**: なし（並列化可能）
**検証**: `gqlgen generate`が成功する

- `graph/schema.graphqls`を編集
  - `User`型定義（id, name, email, createdAt）
  - `Query`に`users: [User!]!`を追加
  - `Query`に`user(id: ID!): User`を追加

### 8. gqlgenコード生成
**目的**: スキーマから型とResolverスケルトンを生成する
**依存**: タスク7（スキーマ拡張）
**検証**: 生成されたコードがコンパイル可能

- `gqlgen generate`を実行
- 生成されたファイルを確認
  - `graph/model/models_gen.go`に`User`型
  - `graph/schema.resolvers.go`にResolverメソッドスケルトン

### 9. Resolver依存性注入更新
**目的**: ResolverにUserRepositoryを注入する
**依存**: タスク5（UserRepositoryインターフェース）
**検証**: `graph/resolver.go`がコンパイル可能

- `graph/resolver.go`を編集
  - `Resolver`構造体に`userRepo repository.UserRepository`フィールド追加
  - `NewResolver`に`userRepo`引数追加

### 10. User Resolverメソッド実装
**目的**: GraphQLクエリからリポジトリを呼び出す
**依存**: タスク6（Repository実装）、タスク9（依存性注入）
**検証**: GraphQL Playgroundでクエリ実行可能

- `graph/schema.resolvers.go`を編集
  - `Users(ctx)`実装: `r.userRepo.List(ctx)`を呼び出し
  - `User(ctx, id)`実装: `r.userRepo.GetByID(ctx, id)`を呼び出し
  - エラーハンドリング
  - `time.Time`から`String`への変換（ISO8601フォーマット）

### 11. server.go統合
**目的**: アプリケーション起動時に全コンポーネントを初期化する
**依存**: タスク4（PostgreSQLクライアント）、タスク9（Resolver更新）
**検証**: `go run server.go`でアプリケーションが起動する

- `server.go`を編集
  - `DATABASE_URL`環境変数を取得（またはデフォルト値）
  - `postgres.NewClient()`でPostgreSQL接続
  - `repository.NewPostgresUserRepository(pgClient)`でリポジトリ初期化
  - `graph.NewResolver(messageRepo, userRepo)`でResolver初期化（両方のリポジトリを渡す）
  - エラーハンドリング（PostgreSQL接続失敗時はgraceful終了）

## Phase 4: Data Seeding & Testing (データ投入とテスト)

### 12. シードスクリプト作成
**目的**: サンプルユーザーデータを投入する
**依存**: タスク2（usersテーブル作成）
**検証**: `go run scripts/seed-postgres.go`でデータ投入成功

- `scripts/seed-postgres.go`を作成
  - `DATABASE_URL`から接続文字列を取得
  - PostgreSQL接続
  - 3-5人のサンプルユーザーをINSERT
  - べき等性: `INSERT ... ON CONFLICT (id) DO UPDATE`
  - 成功メッセージ出力

### 13. Docker Composeでの統合テスト
**目的**: 全コンポーネントが連携動作することを確認する
**依存**: タスク11（server.go統合）
**検証**: `docker compose up`で全サービスが起動し、GraphQL Playgroundでクエリ実行可能

- `docker-compose.yml`の`app`サービスを更新
  - `environment`に`DATABASE_URL`を追加
  - `depends_on`に`postgres`を追加
- `docker compose up`を実行
- シードスクリプトを手動実行: `docker compose exec app go run scripts/seed-postgres.go`

### 14. GraphQL Playgroundでの手動テスト
**目的**: エンドツーエンドで機能確認する
**依存**: タスク13（統合テスト）
**検証**: すべてのクエリが期待通りのレスポンスを返す

- http://localhost:8080/ にアクセス
- `{ users { id name email createdAt } }`クエリを実行
- `{ user(id: "user1") { id name email createdAt } }`クエリを実行
- 存在しないID（`user(id: "nonexistent")`）でエラーレスポンス確認
- 既存の`messages`クエリが引き続き動作することを確認（後方互換性）

## Phase 5: Documentation (ドキュメンテーション)

### 15. README更新
**目的**: ユーザーがPostgreSQL機能を使えるよう手順を記載する
**依存**: タスク14（手動テスト完了）
**検証**: READMEに従ってセットアップ可能

- `README.md`を更新
  - PostgreSQLの起動手順
  - シードスクリプトの実行方法
  - GraphQL Playgroundでのクエリ例
  - 環境変数の説明（`DATABASE_URL`）
  - トラブルシューティング（接続エラー時の対処）

---

## 並列化可能なタスク

以下のタスクは並列実行可能です:

- **グループA（Docker/Infrastructure）**: タスク1, 2
- **グループB（Domain/Repository）**: タスク3, 5
- **グループC（GraphQL Schema）**: タスク7

タスク4は、タスク1完了後に開始できます。
タスク6は、タスク4と5の完了後に開始できます。
タスク8は、タスク7完了後に開始できます。
タスク9以降は順次実行が推奨されます。

## 検証ポイント

各タスク完了後の検証項目:
- **コンパイル**: 変更後に`go build`が成功する
- **フォーマット**: `gofmt -w .`と`goimports -w .`を実行
- **gqlgen**: スキーマ変更後は必ず`gqlgen generate`を実行
- **Docker Compose**: サービス変更後は`docker compose up`で起動確認
- **ログ確認**: エラーログが出ていないか確認

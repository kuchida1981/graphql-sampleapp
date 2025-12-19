# user-repository Specification Delta

## ADDED Requirements

### Requirement: Repository Pattern Implementation
システムはリポジトリパターンを使用してPostgreSQLデータアクセスを抽象化することをMUSTとする。

#### Scenario: UserRepositoryインターフェース定義
- **WHEN** リポジトリコードを確認する
- **THEN** `UserRepository` インターフェースが定義されている
- **AND** `GetByID(ctx context.Context, id string) (*User, error)` メソッドが含まれる
- **AND** `List(ctx context.Context) ([]*User, error)` メソッドが含まれる

#### Scenario: PostgresUserRepository実装
- **WHEN** PostgreSQL実装を確認する
- **THEN** `PostgresUserRepository` 構造体が `UserRepository` を実装している
- **AND** `sql.DB`をフィールドとして持つ
- **AND** コンストラクタで`sql.DB`を受け取る

#### Scenario: Dependency Injection
- **WHEN** Resolverの初期化コードを確認する
- **THEN** `NewResolver` コンストラクタが `MessageRepository` と `UserRepository` を引数として受け取る
- **AND** Resolverが両方のリポジトリをフィールドとして保持する
- **AND** `server.go` でリポジトリの具象実装をResolverに注入する

### Requirement: User Data Operations
システムは`database/sql`を使用してユーザーデータの読み取りを実行することをMUSTとする。

#### Scenario: 全ユーザーの取得
- **WHEN** `List()` メソッドが呼び出される
- **THEN** `SELECT id, name, email, created_at FROM users ORDER BY created_at DESC` クエリを実行する
- **AND** 各行をUser構造体にマッピングする
- **AND** User配列を返す

#### Scenario: 単一ユーザーの取得
- **WHEN** `GetByID("1")` メソッドが呼び出される
- **THEN** `SELECT id, name, email, created_at FROM users WHERE id = $1` クエリを実行する
- **AND** ユーザーが存在すればUser構造体にマッピングする
- **AND** 存在しなければ`sql.ErrNoRows`をハンドリングし、エラーを返す

#### Scenario: プリペアドステートメントの使用
- **WHEN** GetByIDでパラメータを使用する
- **THEN** SQLインジェクション対策としてプレースホルダー (`$1`) を使用する
- **AND** `QueryRowContext(ctx, query, id)` で安全に実行する

#### Scenario: データ型のマッピング
- **WHEN** PostgreSQLの行をGoの構造体にマッピングする
- **THEN** `id`, `name`, `email` は文字列型にマッピングする
- **AND** `created_at` は `time.Time` 型にマッピングする
- **AND** `Scan()`メソッドで各フィールドを読み込む

### Requirement: Error Handling
システムはPostgreSQL操作のエラーを適切にハンドリングすることをMUSTとする。

#### Scenario: ユーザー未存在エラー
- **WHEN** 存在しないIDで `GetByID("nonexistent")` が呼び出される
- **THEN** `sql.ErrNoRows` がキャッチされる
- **AND** `nil` とエラーメッセージ "user not found" を返す

#### Scenario: データベース接続エラー
- **WHEN** クエリ実行中にデータベース接続が失敗する
- **THEN** エラーがログに記録される
- **AND** エラーが呼び出し元に返される

#### Scenario: スキャンエラー
- **WHEN** `Scan()`でデータ型のマッピングに失敗する
- **THEN** エラーがログに記録される
- **AND** エラーが呼び出し元に返される

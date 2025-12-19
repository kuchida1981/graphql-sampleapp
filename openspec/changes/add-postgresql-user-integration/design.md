# Design: PostgreSQL User Integration

## Architecture Overview

この変更は、既存のFirestore Message統合と同じアーキテクチャパターンを踏襲し、PostgreSQLバックエンドでUserドメインを追加する。複数データベースの統合パターンを学習できるよう、意図的に並行した構造を採用する。

```
┌─────────────────────────────────────────────────────────┐
│                   GraphQL Layer                         │
│  (schema.graphqls, Resolver)                           │
└───────────────┬─────────────────┬───────────────────────┘
                │                 │
                │                 │
    ┌───────────▼──────┐  ┌──────▼────────────┐
    │ MessageRepository│  │  UserRepository   │
    │   (interface)    │  │   (interface)     │
    └───────────┬──────┘  └──────┬────────────┘
                │                 │
        ┌───────▼────────┐ ┌─────▼──────────┐
        │   Firestore    │ │   PostgreSQL   │
        │  (Message実装)  │ │   (User実装)   │
        └────────────────┘ └────────────────┘
```

## Component Design

### 1. Domain Layer (`internal/domain`)
**新規追加**: `user.go`
```go
type User struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
}
```

**理由**: Messageと同様のシンプルな構造。PostgreSQLのカラムと1:1でマッピング。

### 2. Repository Layer (`internal/repository`)

**新規追加**: `user.go` (インターフェース定義)
```go
type UserRepository interface {
    List(ctx context.Context) ([]*domain.User, error)
    GetByID(ctx context.Context, id string) (*domain.User, error)
}
```

**新規追加**: `postgres_user.go` (PostgreSQL実装)
```go
type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository
func (r *PostgresUserRepository) List(ctx context.Context) ([]*domain.User, error)
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error)
```

**設計判断**:
- `database/sql`を使用: 標準的で教育的。pgxの高度な機能は不要。
- エラーハンドリング: `sql.ErrNoRows`を適切にハンドリング（GetByIDでnilとエラー返却）
- クエリ: プリペアドステートメントを使用してSQLインジェクション対策

### 3. PostgreSQL Package (`internal/postgres`)

**新規追加**: `client.go`
```go
func NewClient(ctx context.Context, connStr string) (*sql.DB, error)
```

**責務**:
- PostgreSQL接続文字列の構築
- `sql.DB`の初期化とping確認
- 接続プールの設定（デフォルトまたは環境変数から）

**設計判断**:
- Firestoreパッケージ（`internal/firestore`）と並行した構造
- 接続文字列は環境変数から取得（`DATABASE_URL`または個別のパラメータ）
- エラー時はログ出力してエラー返却（graceful degradation）

### 4. GraphQL Layer

**変更**: `graph/schema.graphqls`
```graphql
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: String!
}

extend type Query {
  users: [User!]!
  user(id: ID!): User
}
```

**変更**: `graph/resolver.go`
```go
type Resolver struct {
    messageRepo repository.MessageRepository
    userRepo    repository.UserRepository  // 追加
}

func NewResolver(
    messageRepo repository.MessageRepository,
    userRepo repository.UserRepository,  // 追加
) *Resolver
```

**設計判断**:
- 既存のResolverに`userRepo`フィールドを追加
- 依存性注入で柔軟性を確保
- Firestoreとは独立したリポジトリとして扱う

### 5. Main Entrypoint

**変更**: `server.go`
```go
// PostgreSQLクライアント初期化
pgClient, err := postgres.NewClient(ctx, getDatabaseURL())

// リポジトリ初期化
userRepo := repository.NewPostgresUserRepository(pgClient)

// Resolver初期化（両方のリポジトリを注入）
resolver := graph.NewResolver(messageRepo, userRepo)
```

**設計判断**:
- PostgreSQL接続失敗時はアプリケーションをgracefulに終了（Firestoreと同じ）
- リポジトリはResolverに明示的に注入
- 接続文字列は環境変数から取得（`DATABASE_URL`または`POSTGRES_*`）

## Database Schema

**新規追加**: `scripts/init-postgres.sql`
```sql
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

**設計判断**:
- シンプルなテーブル定義（正規化不要）
- メールアドレスにUNIQUE制約
- 検索性能のためインデックス追加
- タイムスタンプはPostgreSQLのTIMESTAMP型

## Seeding Strategy

**新規追加**: `scripts/seed-postgres.go`
```go
func main() {
    // DATABASE_URL取得
    // PostgreSQL接続
    // サンプルユーザー3-5人をINSERT
    // べき等性: UPSERT (ON CONFLICT DO UPDATE) またはDELETE + INSERT
}
```

**実行方法**:
```bash
go run scripts/seed-postgres.go
```

**設計判断**:
- Firestoreのシードスクリプトと並行した構造
- べき等性を確保（複数回実行可能）
- 環境変数で接続先を切り替え可能

## Docker Compose Integration

**変更**: `docker-compose.yml`
```yaml
services:
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: graphql_user
      POSTGRES_PASSWORD: graphql_pass
      POSTGRES_DB: graphql_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-postgres.sql:/docker-entrypoint-initdb.d/init.sql

  app:
    # ... 既存設定 ...
    environment:
      # ... 既存環境変数 ...
      - DATABASE_URL=postgres://graphql_user:graphql_pass@postgres:5432/graphql_db?sslmode=disable
    depends_on:
      - firestore
      - postgres  # 追加

volumes:
  postgres_data:
```

**設計判断**:
- PostgreSQL 16 Alpine: 軽量で最新の安定版
- 初期化SQLスクリプトを`docker-entrypoint-initdb.d`に配置（自動実行）
- ボリュームでデータ永続化
- アプリケーションはPostgresの起動を待つ（`depends_on`）

## Error Handling

### PostgreSQL接続エラー
- `postgres.NewClient()`でエラー発生 → ログ出力 + アプリケーション終了
- Firestoreと同じパターンで一貫性確保

### クエリエラー
- `sql.ErrNoRows` (GetByID) → `nil, error`を返し、ResolverでGraphQLエラー化
- その他のDBエラー → ログ出力 + エラー返却

### ロギング
- Firestoreと同じログフォーマット
- 操作開始・完了・エラーを記録

## Testing Strategy

### 統合テスト（将来的）
- テスト用PostgreSQLコンテナを起動
- リポジトリの各メソッドをテスト
- Firestoreテストと並行した構造

### 手動テスト
- GraphQL Playgroundで`users`/`user`クエリを実行
- シードスクリプトでデータ投入後、取得確認

## Trade-offs and Alternatives

### 選択: `database/sql` vs `pgx`ネイティブAPI
- **選択**: `database/sql`
- **理由**: 標準的で学習曲線が緩やか。このチュートリアルの目的に合致。
- **代替案**: `pgx`のネイティブAPIは高性能だが、複雑性が増す。

### 選択: 初期化SQLスクリプト vs マイグレーションツール
- **選択**: 初期化SQLスクリプト
- **理由**: シンプルで理解しやすい。本番環境でないため十分。
- **代替案**: golang-migrate/gooseは本番的だが、学習コストが高い。

### 選択: UserとMessageの独立 vs 関連づけ
- **選択**: 完全に独立したドメインとして扱う
- **理由**: スコープを絞り、複雑性を避ける。将来的に関連づけ可能。
- **代替案**: User.MessagesやMessage.Authorを実装すると、リレーショナルDBの利点を示せるが、今回はスコープ外。

## Dependencies

### 新規依存関係
- `github.com/jackc/pgx/v5/stdlib`: pgxのdatabase/sqlドライバ
- `github.com/lib/pq`: 代替PostgreSQLドライバ（または）

**選択**: `pgx/v5/stdlib`
**理由**: 現代的で活発にメンテナンスされている。将来的にネイティブAPIに移行しやすい。

### 環境変数
- `DATABASE_URL`: PostgreSQL接続文字列（優先）
- または `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`

## Future Considerations

この設計は将来的な拡張を考慮している:
1. **UserとMessageの関連**: `Message.authorID`でUserを参照
2. **DataLoader**: N+1問題回避のためのバッチローディング
3. **Mutation**: Create/Update/Delete操作の追加
4. **ページネーション**: カーソルベースまたはオフセットベース
5. **トランザクション**: 複数テーブルの更新

これらは明示的にスコープ外としているが、今回の設計はこれらの追加を妨げない。
# Design: Firestore Docker Integration

## Architecture Overview
このデザインでは、既存のGraphQLサーバーにCloud Firestoreバックエンドを追加し、Docker Composeで管理されるマルチコンテナ環境を構築します。リポジトリパターンを採用し、データアクセスロジックをリゾルバーから分離することで、テスタビリティと保守性を確保します。

## System Components

### 1. Docker Compose Architecture
```
┌─────────────────────────────────────────┐
│         Docker Compose Network          │
│                                         │
│  ┌──────────────┐    ┌──────────────┐  │
│  │     app      │───▶│  firestore   │  │
│  │ (Go/GraphQL) │    │  (Emulator)  │  │
│  │   :8080      │    │   :8081      │  │
│  └──────────────┘    └──────────────┘  │
│         │                               │
└─────────┼───────────────────────────────┘
          │
    ┌─────▼─────┐
    │  Client   │
    │ (Browser) │
    └───────────┘
```

**コンポーネント説明:**
- **app**: GraphQLサーバーを実行するGoアプリケーションコンテナ
- **firestore**: Google Cloud Firestore Emulator（ローカル開発用）
- **ネットワーク**: Docker Composeのデフォルトブリッジネットワークでコンテナ間通信

### 2. Application Layer Architecture
```
GraphQL Layer (Handler/Playground)
         │
         ▼
    Resolvers
         │
         ▼
   Service Layer (Future)
         │
         ▼
  Repository Interface
         │
         ▼
Firestore Repository Impl
         │
         ▼
   Firestore Client
         │
         ▼
  Firestore Emulator
```

**レイヤー責務:**
- **Resolvers**: GraphQLリクエストを受け取り、リポジトリを呼び出す
- **Repository Interface**: データアクセスの抽象化（将来的にモック可能）
- **Firestore Repository**: Firestore固有の実装
- **Firestore Client**: Firebase Admin SDKによるFirestore接続

## Key Design Decisions

### Decision 1: Firestore Emulator vs Real Firestore
**選択**: ローカル開発にはFirestore Emulatorを使用

**理由**:
- GCPプロジェクトや認証情報が不要
- 完全にローカルで動作し、インターネット接続不要
- 高速なイテレーション（データリセットが容易）
- 開発コストがゼロ

**トレードオフ**:
- Emulatorは実際のFirestoreと100%同一ではない（一部機能制限）
- 本番環境との設定切り替えが必要

### Decision 2: Repository Pattern
**選択**: Repositoryパターンでデータアクセスを抽象化

**理由**:
- プロジェクトの規約（project.md: "Repositoryパターン: リポジトリインターフェースの背後にデータベースアクセスを抽象化"）に準拠
- テスタビリティ向上（モックによるユニットテスト）
- 将来的にPostgreSQLなど他のデータソース追加時に変更箇所を局所化

**実装方針**:
```go
// Repository interface
type MessageRepository interface {
    GetByID(ctx context.Context, id string) (*Message, error)
    List(ctx context.Context) ([]*Message, error)
}

// Firestore implementation
type FirestoreMessageRepository struct {
    client *firestore.Client
}
```

### Decision 3: Dependency Injection
**選択**: コンストラクタインジェクションでリポジトリをリゾルバーに渡す

**理由**:
- プロジェクト規約に準拠（"依存性注入: コンストラクタを通じて依存関係を明示的に渡す"）
- テスト時にモックリポジトリを注入可能
- 依存関係が明確

**実装方針**:
```go
type Resolver struct {
    messageRepo MessageRepository
}

func NewResolver(messageRepo MessageRepository) *Resolver {
    return &Resolver{messageRepo: messageRepo}
}
```

### Decision 4: Environment-based Configuration
**選択**: 環境変数でFirestoreエンドポイントを切り替え

**理由**:
- 開発環境（Emulator）と本番環境を同一コードで対応
- Docker Composeで環境変数を簡単に設定可能
- 12 Factor Appの原則に準拠

**環境変数**:
- `FIRESTORE_EMULATOR_HOST`: Emulatorのホスト（例: `firestore:8081`）
- `GCP_PROJECT_ID`: プロジェクトID（Emulator使用時はダミー値可）

### Decision 5: Sample Domain Model
**選択**: シンプルな "Message" エンティティをサンプルとして使用

**理由**:
- 学習目的のプロジェクトであり、複雑なドメインモデルは不要
- 単一エンティティで CRUD の基本パターンを示せる
- 将来的に拡張しやすい

**Message モデル**:
```
Message {
  id: String!
  content: String!
  author: String!
  createdAt: String!
}
```

## Data Flow

### Query Flow Example: `messages` クエリ
1. Client → GraphQL Server: `{ messages { id content author } }`
2. GraphQL Handler → Resolver: `Query.messages()`
3. Resolver → Repository: `messageRepo.List(ctx)`
4. Repository → Firestore Client: Firestore API Call
5. Firestore Client → Emulator: gRPC/HTTP Request
6. Emulator → Firestore Client: Document data
7. Repository: Convert Firestore docs to Go structs
8. Resolver: Return to GraphQL layer
9. GraphQL Server → Client: JSON response

## Docker Configuration

### Dockerfile Strategy
**選択**: マルチステージビルドを使用

**理由**:
- ビルド成果物のみを含む最小限のイメージサイズ
- Go の静的バイナリ特性を活用
- セキュリティ向上（不要なビルドツールを含まない）

**構成**:
```dockerfile
# Stage 1: Build
FROM golang:1.21 AS builder
# ... build steps

# Stage 2: Runtime
FROM alpine:latest
# ... copy binary and run
```

### Docker Compose Strategy
**選択**: サービス間でネットワークを共有し、サービス名でDNS解決

**理由**:
- シンプルな設定
- Docker標準のサービスディスカバリー機能を活用
- `FIRESTORE_EMULATOR_HOST=firestore:8081` のような直感的な設定

## Error Handling
- Firestore接続エラー: アプリケーション起動時にログ出力し、gracefulに処理
- ドキュメント未存在: GraphQLの標準エラーレスポンスとして返却
- Emulator未起動: 接続リトライロジックなし（明示的なエラーメッセージで開発者に通知）

## Testing Strategy
- **Unit Tests**: リポジトリインターフェースのモックを使用してリゾルバーをテスト
- **Integration Tests**: Firestore Emulatorを使用した結合テスト（将来的に追加）
- **Manual Testing**: Docker Compose環境でGraphQL Playgroundを使用

## Migration Path
現在のHello World実装からの移行は非破壊的:
1. 既存の `hello` クエリは維持
2. 新しい `messages` クエリを追加
3. Docker Composeはオプション（`go run server.go` も引き続き動作）

## Future Considerations
- PostgreSQL追加時: 別のリポジトリ実装を追加（同じパターン）
- Mutation追加時: Repository interfaceにCreate/Update/Deleteメソッドを追加
- 認証追加時: Context経由でユーザー情報を伝播
- プロダクションデプロイ: 環境変数で実際のFirestoreプロジェクトに接続
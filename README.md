# GraphQL Sample App

GraphQLのベストプラクティスとパターンを学習するためのチュートリアルプロジェクト。gqlgenを使用したプロダクション対応GraphQL APIの構築方法を示すリファレンス実装。

## 技術スタック

- **言語:** Go (Golang)
- **GraphQLフレームワーク:** gqlgen
- **データベース:**
  - PostgreSQL - リレーショナルデータ（ユーザー、気象アラートメタデータ）
  - Cloud Firestore - ドキュメントストア（メッセージ、気象アラート詳細データ）
- **開発環境:** Docker Compose

## はじめに

### 前提条件

- Go 1.24.0以上
- Docker および Docker Compose (推奨)

### Docker Composeでの起動 (推奨)

Docker Composeを使用すると、Firestore Emulatorを含む完全な開発環境を簡単に起動できます。

```bash
# 環境変数ファイルを作成（オプション）
cp .env.example .env

# すべてのサービスを起動
docker-compose up

# バックグラウンドで起動
docker-compose up -d

# ログを確認
docker-compose logs -f app
```

サーバーが起動すると、`http://localhost:8080/` でGraphQL Playgroundにアクセスできます。

#### サンプルデータの投入

Docker Composeで起動している場合、コンテナ内で実行します:

```bash
# PostgreSQLにサンプルユーザーを投入
docker-compose exec app go run scripts/seed-postgres.go

# Firestore Emulatorにサンプルメッセージを投入
docker-compose exec app go run scripts/seed-firestore.go

# PostgreSQL + Firestoreにサンプル気象アラートを投入
docker-compose exec app go run scripts/seed-weather-alerts.go
```

ローカル環境で実行する場合:

```bash
# PostgreSQLにサンプルユーザーを投入
export DATABASE_URL="postgres://graphql_user:graphql_pass@localhost:5432/graphql_db?sslmode=disable"
go run scripts/seed-postgres.go

# Firestore Emulatorにサンプルメッセージを投入
export FIRESTORE_EMULATOR_HOST=localhost:8081
export GCP_PROJECT_ID=demo-project
go run scripts/seed-firestore.go

# PostgreSQL + Firestoreにサンプル気象アラートを投入
export DATABASE_URL="postgres://graphql_user:graphql_pass@localhost:5432/graphql_db?sslmode=disable"
export FIRESTORE_EMULATOR_HOST=localhost:8081
export GCP_PROJECT_ID=demo-project
go run scripts/seed-weather-alerts.go
```

#### 環境の停止とクリーンアップ

```bash
# サービスを停止
docker-compose down

# データを含めて完全にクリーンアップ
docker-compose down -v
```

### ローカル環境での起動

Docker Composeを使用しない場合は、以下の手順でローカル環境で起動できます。

#### インストール

```bash
# リポジトリのクローン
git clone <repository-url>
cd graphql-sampleapp

# 依存関係のインストール
go mod download
```

### サーバーの起動

```bash
go run server.go
```

サーバーが起動すると、以下のメッセージが表示されます:

```
connect to http://localhost:8080/ for GraphQL playground
```

## 使い方

### GraphQL Playground

ブラウザで http://localhost:8080/ にアクセスすると、GraphQL Playgroundが開きます。

### クエリ例

#### Hello Worldクエリ

```graphql
{
  hello
}
```

レスポンス:

```json
{
  "data": {
    "hello": "Hello World"
  }
}
```

#### 全メッセージの取得

```graphql
{
  messages {
    id
    content
    author
    createdAt
  }
}
```

レスポンス:

```json
{
  "data": {
    "messages": [
      {
        "id": "msg3",
        "content": "Docker Compose makes local development easy.",
        "author": "Charlie",
        "createdAt": "2025-12-19T10:00:00+09:00"
      },
      {
        "id": "msg2",
        "content": "GraphQL and Firestore integration is working!",
        "author": "Bob",
        "createdAt": "2025-12-19T09:00:00+09:00"
      },
      {
        "id": "msg1",
        "content": "Hello, Firestore! This is the first message.",
        "author": "Alice",
        "createdAt": "2025-12-19T08:00:00+09:00"
      }
    ]
  }
}
```

#### 特定のメッセージを取得

```graphql
{
  message(id: "msg1") {
    id
    content
    author
    createdAt
  }
}
```

レスポンス:

```json
{
  "data": {
    "message": {
      "id": "msg1",
      "content": "Hello, Firestore! This is the first message.",
      "author": "Alice",
      "createdAt": "2025-12-19T08:00:00+09:00"
    }
  }
}
```

#### 全ユーザーの取得

```graphql
{
  users {
    id
    name
    email
    createdAt
  }
}
```

レスポンス:

```json
{
  "data": {
    "users": [
      {
        "id": "user5",
        "name": "Eve Adams",
        "email": "eve@example.com",
        "createdAt": "2025-12-19T17:00:00+09:00"
      },
      {
        "id": "user4",
        "name": "Diana Prince",
        "email": "diana@example.com",
        "createdAt": "2025-12-19T11:00:00+09:00"
      }
    ]
  }
}
```

#### 特定のユーザーを取得

```graphql
{
  user(id: "user1") {
    id
    name
    email
    createdAt
  }
}
```

レスポンス:

```json
{
  "data": {
    "user": {
      "id": "user1",
      "name": "Alice Smith",
      "email": "alice@example.com",
      "createdAt": "2025-12-17T17:00:00+09:00"
    }
  }
}
```

#### 全気象アラートの取得（ハイブリッドクエリ）

PostgreSQLでメタデータを検索し、Firestoreから詳細データを取得するハイブリッドパターンのデモ:

```graphql
{
  weatherAlerts {
    id
    region
    severity
    issuedAt
    title
    description
    rawData
    affectedAreas
    recommendations
  }
}
```

#### 地域で気象アラートをフィルタ

```graphql
{
  weatherAlerts(region: "Tokyo") {
    id
    region
    severity
    issuedAt
    title
  }
}
```

#### 日時で気象アラートをフィルタ

```graphql
{
  weatherAlerts(issuedAfter: "2025-12-19T00:00:00Z") {
    id
    region
    severity
    issuedAt
    title
  }
}
```

#### データフロー（PostgreSQL → Firestore）

1. **PostgreSQL**: メタデータ検索（地域・重要度・発行日時でフィルタ）
2. **結果**: IDリストを取得
3. **Firestore**: IDリストからバッチで詳細データ取得
4. **マージ**: PostgreSQLメタデータ + Firestore詳細データ
5. **レスポンス**: 統合されたWeatherAlertデータ

このパターンにより、スキーマ変更に強い柔軟なデータ構造（Firestore）と効率的な検索・集計（PostgreSQL）を両立できます。

### cURLでのクエリ実行

```bash
# Hello Worldクエリ
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ hello }"}'

# 全メッセージの取得
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ messages { id content author createdAt } }"}'

# 全ユーザーの取得
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ users { id name email createdAt } }"}'

# 特定のユーザーを取得
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ user(id: \"user1\") { id name email createdAt } }"}'
```

## プロジェクト構造

```
.
├── server.go              # GraphQLサーバーのエントリーポイント
├── gqlgen.yml             # gqlgen設定ファイル
├── Dockerfile             # Goアプリケーション用のDockerイメージ
├── docker-compose.yml     # Docker Compose設定
├── .env.example           # 環境変数のサンプル
├── graph/
│   ├── schema.graphqls    # GraphQLスキーマ定義
│   ├── resolver.go        # リゾルバーのベース構造
│   ├── schema.resolvers.go # リゾルバー実装
│   ├── generated.go       # gqlgenが生成したコード
│   └── model/             # GraphQLモデルの型定義
├── internal/
│   ├── domain/            # ドメインモデル
│   │   ├── message.go     # Messageエンティティ
│   │   └── user.go        # Userエンティティ
│   ├── firestore/         # Firestoreクライアント
│   │   └── client.go      # Firestore初期化
│   ├── postgres/          # PostgreSQLクライアント
│   │   └── client.go      # PostgreSQL初期化
│   └── repository/        # データアクセス層
│       ├── message.go     # MessageRepositoryインターフェース
│       ├── user.go        # UserRepositoryインターフェース
│       ├── firestore_message.go # Firestore Message実装
│       └── postgres_user.go     # PostgreSQL User実装
├── scripts/
│   ├── init-postgres.sql  # PostgreSQL初期化スクリプト
│   ├── seed-postgres.go   # PostgreSQLサンプルデータシード
│   └── seed-firestore.go  # Firestoreサンプルデータシード
├── go.mod                 # Go module定義
└── go.sum                 # Go依存関係のチェックサム
```

## 開発

### スキーマの変更

1. `graph/schema.graphqls` を編集
2. コード生成を実行:

```bash
go run github.com/99designs/gqlgen generate
```

または、gqlgenをインストール済みの場合:

```bash
~/go/bin/gqlgen generate
```

3. 新しいリゾルバーを `graph/schema.resolvers.go` に実装

## ライセンス

MIT

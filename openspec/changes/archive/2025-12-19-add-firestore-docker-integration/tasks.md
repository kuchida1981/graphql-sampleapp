# Tasks: Add Firestore Docker Integration

このタスクリストは、Firestore統合とDocker Compose環境を段階的に実装するための作業項目です。各タスクは小さく検証可能な単位に分割され、順序立てて実装することで安全に機能を追加できます。

## Phase 1: Docker Infrastructure Setup

### Task 1.1: Create Dockerfile for Go application
**目的**: Goアプリケーション用のマルチステージDockerイメージを作成

**作業内容**:
- プロジェクトルートに `Dockerfile` を作成
- ビルドステージ: `golang:1.21` ベースイメージで依存関係インストールとビルド
- ランタイムステージ: `alpine:latest` で最小限のイメージを作成
- ポート8080を公開
- エントリーポイントでバイナリを実行

**検証方法**:
```bash
docker build -t graphql-sampleapp .
docker run -p 8080:8080 graphql-sampleapp
# ブラウザで http://localhost:8080 にアクセスし、GraphQL Playgroundが表示されることを確認
```

**依存関係**: なし

---

### Task 1.2: Create docker-compose.yml with app service
**目的**: Docker Compose設定を作成し、appサービスを定義

**作業内容**:
- プロジェクトルートに `docker-compose.yml` を作成
- `app` サービスを定義（Dockerfileからビルド）
- ポートマッピング `8080:8080` を設定
- 環境変数 `PORT=8080` を設定
- ボリュームマウントでホットリロード対応（`./:/app`, 作業ディレクトリ設定）

**検証方法**:
```bash
docker-compose up app
# http://localhost:8080 でGraphQL Playgroundが表示されることを確認
docker-compose down
```

**依存関係**: Task 1.1

---

### Task 1.3: Add Firestore Emulator service to docker-compose.yml
**目的**: Firestore Emulator用のサービスをDocker Composeに追加

**作業内容**:
- `firestore` サービスを追加（`google/cloud-sdk:emulators` イメージ使用）
- コマンド: `gcloud emulators firestore start --host-port=0.0.0.0:8081`
- ポート公開: `8081:8081`
- ヘルスチェック設定（オプション）

**検証方法**:
```bash
docker-compose up firestore
# ログに "Dev App Server is now running" が表示されることを確認
```

**依存関係**: Task 1.2

---

### Task 1.4: Configure service dependencies and environment variables
**目的**: appサービスがfirestoreサービスに依存するよう設定し、環境変数を追加

**作業内容**:
- `app` サービスに `depends_on: [firestore]` を追加
- `app` サービスの環境変数に以下を追加:
  - `FIRESTORE_EMULATOR_HOST=firestore:8081`
  - `GCP_PROJECT_ID=demo-project`

**検証方法**:
```bash
docker-compose up
# firestoreが先に起動し、その後appが起動することをログで確認
# appのログに環境変数が正しく設定されていることを確認
```

**依存関係**: Task 1.3

---

### Task 1.5: Create .env.example file
**目的**: 環境変数のサンプルファイルを提供

**作業内容**:
- `.env.example` ファイルを作成
- 以下の変数を記載:
  ```
  GCP_PROJECT_ID=demo-project
  FIRESTORE_EMULATOR_HOST=firestore:8081
  PORT=8080
  ```
- `.gitignore` に `.env` を追加（すでに追加済みの場合はスキップ）

**検証方法**:
- ファイルが作成されていることを確認
- READMEで `.env.example` を `.env` にコピーする手順を案内

**依存関係**: Task 1.4

---

## Phase 2: Firestore Integration

### Task 2.1: Add Firebase Admin SDK dependency
**目的**: Firebase Admin SDK for Goをプロジェクトに追加

**作業内容**:
- `go get firebase.google.com/go/v4` を実行
- `go mod tidy` でモジュールを整理

**検証方法**:
```bash
go mod verify
# エラーがないことを確認
```

**依存関係**: なし（Phase 1と並行可能）

---

### Task 2.2: Create Firestore client initialization code
**目的**: Firestoreクライアントを初期化するコードを実装

**作業内容**:
- `internal/firestore/client.go` ファイルを作成（ディレクトリも作成）
- `NewClient(ctx context.Context, projectID string) (*firestore.Client, error)` 関数を実装
- `FIRESTORE_EMULATOR_HOST` 環境変数が設定されている場合は自動的にEmulatorに接続
- エラーハンドリングとログ出力を含める

**検証方法**:
- ユニットテスト作成（環境変数モック）
- または、`main.go` で初期化を試行しログを確認

**依存関係**: Task 2.1, Task 1.4

---

### Task 2.3: Define Message domain model
**目的**: Messageエンティティの構造体を定義

**作業内容**:
- `internal/domain/message.go` ファイルを作成
- `Message` 構造体を定義:
  ```go
  type Message struct {
      ID        string
      Content   string
      Author    string
      CreatedAt time.Time
  }
  ```
- Firestoreタグを追加（例: `firestore:"content"`）

**検証方法**:
- コードがコンパイルされることを確認

**依存関係**: なし

---

### Task 2.4: Create MessageRepository interface
**目的**: リポジトリパターンのインターフェースを定義

**作業内容**:
- `internal/repository/message.go` ファイルを作成
- `MessageRepository` インターフェースを定義:
  ```go
  type MessageRepository interface {
      List(ctx context.Context) ([]*domain.Message, error)
      GetByID(ctx context.Context, id string) (*domain.Message, error)
  }
  ```

**検証方法**:
- コードがコンパイルされることを確認

**依存関係**: Task 2.3

---

### Task 2.5: Implement FirestoreMessageRepository
**目的**: Firestore用のリポジトリ実装を作成

**作業内容**:
- `internal/repository/firestore_message.go` ファイルを作成
- `FirestoreMessageRepository` 構造体を実装
- `NewFirestoreMessageRepository(client *firestore.Client) *FirestoreMessageRepository` コンストラクタ
- `List(ctx)` メソッド: `messages` コレクション全体を取得、`createdAt` で降順ソート
- `GetByID(ctx, id)` メソッド: 指定IDのドキュメントを取得
- Firestoreドキュメント ↔ Message構造体のマッピングロジック
- エラーハンドリング（ドキュメント未存在、接続エラーなど）

**検証方法**:
- 統合テスト作成（Firestore Emulator使用）
- または、後続タスクでGraphQL経由でテスト

**依存関係**: Task 2.2, Task 2.4

---

### Task 2.6: Update GraphQL schema with Message type and queries
**目的**: GraphQLスキーマにMessageタイプとクエリを追加

**作業内容**:
- `graph/schema.graphqls` を更新
- `Message` タイプを追加:
  ```graphql
  type Message {
      id: ID!
      content: String!
      author: String!
      createdAt: String!
  }
  ```
- `Query` タイプに以下を追加:
  ```graphql
  messages: [Message!]!
  message(id: ID!): Message
  ```
- `go run github.com/99designs/gqlgen generate` で再生成

**検証方法**:
```bash
go run github.com/99designs/gqlgen generate
# エラーなくコード生成されることを確認
```

**依存関係**: Task 2.3

---

### Task 2.7: Update Resolver to include MessageRepository
**目的**: ResolverにMessageRepositoryを依存性注入

**作業内容**:
- `graph/resolver.go` の `Resolver` 構造体にフィールド追加:
  ```go
  type Resolver struct {
      messageRepo repository.MessageRepository
  }
  ```
- コンストラクタ `NewResolver(messageRepo repository.MessageRepository)` を作成または更新

**検証方法**:
- コードがコンパイルされることを確認

**依存関係**: Task 2.4

---

### Task 2.8: Implement messages and message resolvers
**目的**: GraphQLクエリのリゾルバーを実装

**作業内容**:
- `graph/schema.resolvers.go` を更新（gqlgen生成後に編集）
- `Messages(ctx)` リゾルバー: `r.messageRepo.List(ctx)` を呼び出し結果を返す
- `Message(ctx, id)` リゾルバー: `r.messageRepo.GetByID(ctx, id)` を呼び出し結果を返す
- エラーハンドリング: リポジトリエラーをGraphQLエラーとして返す
- ログ出力（デバッグ用）

**検証方法**:
- 後続タスクでGraphQL Playground経由でテスト

**依存関係**: Task 2.5, Task 2.6, Task 2.7

---

### Task 2.9: Update main.go to initialize Firestore and Repository
**目的**: main関数でFirestoreクライアントとリポジトリを初期化

**作業内容**:
- `server.go` (または `main.go`) を更新
- 環境変数 `GCP_PROJECT_ID` と `FIRESTORE_EMULATOR_HOST` を読み込み
- Firestoreクライアント初期化: `firestore.NewClient(ctx, projectID)`
- `FirestoreMessageRepository` を作成
- Resolverにリポジトリを注入: `graph.NewResolver(messageRepo)`
- エラーハンドリング: 初期化失敗時はログ出力して終了

**検証方法**:
```bash
docker-compose up
# サーバーが正常起動し、Firestore接続成功のログが表示されることを確認
```

**依存関係**: Task 2.2, Task 2.5, Task 2.7

---

## Phase 3: Sample Data and Testing

### Task 3.1: Create Firestore seed script
**目的**: サンプルデータを投入するスクリプトを作成

**作業内容**:
- `scripts/seed-firestore.go` ファイルを作成
- 環境変数 `FIRESTORE_EMULATOR_HOST`, `GCP_PROJECT_ID` を読み込み
- Firestoreクライアント初期化
- 3〜5件のサンプルMessageをFirestoreに書き込み:
  - id: `msg1`, `msg2`, `msg3`
  - content, author, createdAt をダミーデータで設定
- 成功メッセージをコンソール出力

**検証方法**:
```bash
docker-compose up -d firestore
export FIRESTORE_EMULATOR_HOST=localhost:8081
export GCP_PROJECT_ID=demo-project
go run scripts/seed-firestore.go
# "Successfully seeded N messages" が表示されることを確認
```

**依存関係**: Task 2.2, Task 2.3

---

### Task 3.2: Test messages query in GraphQL Playground
**目的**: GraphQL Playgroundでmessagesクエリをテスト

**作業内容**:
1. Docker Compose環境を起動: `docker-compose up`
2. シードスクリプト実行（コンテナ内またはホストから）
3. ブラウザで `http://localhost:8080` にアクセス
4. GraphQL Playgroundで以下を実行:
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
5. レスポンスにサンプルデータが含まれることを確認

**検証方法**:
- クエリが成功しデータが返されることを確認
- エラーがないことを確認

**依存関係**: Task 2.8, Task 2.9, Task 3.1

---

### Task 3.3: Test message(id) query in GraphQL Playground
**目的**: ID指定のmessageクエリをテスト

**作業内容**:
1. GraphQL Playgroundで以下を実行:
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
2. 指定したIDのMessageが返されることを確認
3. 存在しないIDでエラーが返されることを確認:
   ```graphql
   {
       message(id: "nonexistent") {
           id
       }
   }
   ```

**検証方法**:
- 有効なIDで正しいデータが返される
- 無効なIDでエラーメッセージが返される

**依存関係**: Task 3.2

---

## Phase 4: Documentation and Cleanup

### Task 4.1: Update README with Docker Compose instructions
**目的**: READMEにDocker Compose使用方法を追記

**作業内容**:
- `README.md` に新しいセクション "Docker Compose Setup" を追加
- 以下の手順を記載:
  1. 前提条件（Docker/Docker Compose インストール）
  2. 環境変数設定（`.env.example` → `.env` のコピー）
  3. `docker-compose up` で起動
  4. シードスクリプト実行方法
  5. GraphQL Playgroundアクセス方法（`http://localhost:8080`）
  6. 停止方法（`docker-compose down`）
- スクリーンショットまたは例示クエリを含める

**検証方法**:
- READMEの手順に従って第三者が環境をセットアップできることを確認

**依存関係**: Task 3.3

---

### Task 4.2: Add code comments and documentation
**目的**: コードに適切なコメントとドキュメントを追加

**作業内容**:
- 各関数にGoDocコメントを追加
- 複雑なロジック（例: Firestoreマッピング）に説明コメントを追加
- READMEにアーキテクチャ概要図を追加（オプション）

**検証方法**:
- `go doc` でドキュメントが表示されることを確認
- コードレビュー時に読みやすいことを確認

**依存関係**: 全Phase 1-3タスク完了後

---

### Task 4.3: Run gofmt and goimports
**目的**: コードフォーマットをプロジェクト規約に準拠

**作業内容**:
- `gofmt -w .` を実行
- `goimports -w .` を実行（未インストールの場合は `go install golang.org/x/tools/cmd/goimports@latest`）
- インポートが3セクション（標準、外部、内部）に分かれていることを確認

**検証方法**:
```bash
gofmt -l .
# 出力がないことを確認（全ファイルがフォーマット済み）
```

**依存関係**: Task 4.2

---

### Task 4.4: Final integration test
**目的**: 全機能が統合されて動作することを最終確認

**作業内容**:
1. クリーン状態から開始: `docker-compose down -v`
2. 起動: `docker-compose up -d`
3. シードスクリプト実行
4. GraphQL Playgroundで以下をテスト:
   - `hello` クエリ（既存機能）
   - `messages` クエリ（新機能）
   - `message(id)` クエリ（新機能）
5. ログにエラーがないことを確認
6. 環境停止: `docker-compose down`

**検証方法**:
- 全クエリが成功する
- エラーログがない
- READMEの手順が正確

**依存関係**: Task 4.3

---

## Parallelizable Tasks
以下のタスクは並行して作業可能:
- **Phase 1全体** と **Task 2.1, 2.3** は並行可能
- **Task 2.3** と **Task 2.4** は順序依存だが、**Task 2.1, 2.2** とは独立
- **Task 3.1** と **Task 3.2** の間でシードを実行するタイミング調整が必要

## Critical Path
最も時間がかかる可能性のあるパス:
1. Task 1.1 → 1.2 → 1.3 → 1.4 → Task 2.9 → Task 3.2 → Task 3.3 → Task 4.4
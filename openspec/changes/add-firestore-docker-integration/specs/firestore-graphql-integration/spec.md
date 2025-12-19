# firestore-graphql-integration Specification

## Purpose
Cloud FirestoreをGraphQL APIのバックエンドとして統合し、Firestoreからデータを取得して返却するサンプル機能を提供する。リポジトリパターンを使用してデータアクセスを抽象化し、テスタブルで保守しやすいアーキテクチャを実現する。

## Dependencies
- `docker-compose-setup`: Firestore Emulatorの起動に依存

## ADDED Requirements

### Requirement: Firestore Client Initialization
システムはFirebase Admin SDKを使用してFirestoreクライアントを初期化することをMUSTとする。

#### Scenario: Firestoreクライアントの初期化成功
- **WHEN** アプリケーションが起動する
- **THEN** `FIRESTORE_EMULATOR_HOST` 環境変数が設定されている場合、Emulatorに接続する
- **AND** Firestoreクライアントが正常に初期化される
- **AND** 初期化ログが標準出力に記録される

#### Scenario: Emulator接続の自動検出
- **WHEN** `FIRESTORE_EMULATOR_HOST` 環境変数が設定されている
- **THEN** 認証情報なしでEmulatorに接続する
- **AND** `GCP_PROJECT_ID` のダミー値（例: `demo-project`）を使用する

#### Scenario: 初期化エラーハンドリング
- **WHEN** Firestore初期化が失敗する
- **THEN** エラーログが出力される
- **AND** アプリケーションはgracefulに終了する

### Requirement: Repository Pattern Implementation
システムはリポジトリパターンを使用してFirestoreデータアクセスを抽象化することをMUSTとする。

#### Scenario: MessageRepositoryインターフェース定義
- **WHEN** リポジトリコードを確認する
- **THEN** `MessageRepository` インターフェースが定義されている
- **AND** `GetByID(ctx context.Context, id string) (*Message, error)` メソッドが含まれる
- **AND** `List(ctx context.Context) ([]*Message, error)` メソッドが含まれる

#### Scenario: FirestoreMessageRepository実装
- **WHEN** Firestore実装を確認する
- **THEN** `FirestoreMessageRepository` 構造体が `MessageRepository` を実装している
- **AND** Firestoreクライアントをフィールドとして持つ
- **AND** コンストラクタでFirestoreクライアントを受け取る

#### Scenario: Dependency Injection
- **WHEN** Resolverの初期化コードを確認する
- **THEN** `NewResolver` コンストラクタが `MessageRepository` を引数として受け取る
- **AND** Resolverがリポジトリをフィールドとして保持する
- **AND** `main.go` でリポジトリの具象実装をResolverに注入する

### Requirement: GraphQL Schema Extension
システムはFirestoreからデータを取得するためのGraphQLクエリをスキーマに追加することをMUSTとする。

#### Scenario: Messageタイプの定義
- **WHEN** GraphQLスキーマファイル（`schema.graphqls`）を確認する
- **THEN** `Message` タイプが定義されている
- **AND** `id: ID!` フィールドが含まれる
- **AND** `content: String!` フィールドが含まれる
- **AND** `author: String!` フィールドが含まれる
- **AND** `createdAt: String!` フィールドが含まれる

#### Scenario: messagesクエリの定義
- **WHEN** GraphQLスキーマの `Query` タイプを確認する
- **THEN** `messages: [Message!]!` クエリが定義されている
- **AND** 既存の `hello` クエリは維持されている

#### Scenario: messageクエリの定義（ID指定）
- **WHEN** GraphQLスキーマの `Query` タイプを確認する
- **THEN** `message(id: ID!): Message` クエリが定義されている
- **AND** 引数 `id` が必須（`ID!`）である

### Requirement: Resolver Implementation
システムはFirestoreリポジトリを使用してGraphQLクエリを解決することをMUSTとする。

#### Scenario: messagesリゾルバーの実装
- **WHEN** `messages` クエリが実行される
- **THEN** Resolverは `messageRepo.List(ctx)` を呼び出す
- **AND** Firestoreから全Messageドキュメントを取得する
- **AND** Message配列をGraphQLレスポンスとして返す

#### Scenario: messageリゾルバーの実装（ID指定）
- **WHEN** `message(id: "123")` クエリが実行される
- **THEN** Resolverは `messageRepo.GetByID(ctx, "123")` を呼び出す
- **AND** 指定されたIDのMessageドキュメントを取得する
- **AND** 該当するMessageをGraphQLレスポンスとして返す

#### Scenario: ドキュメント未存在のエラーハンドリング
- **WHEN** 存在しないIDで `message(id: "nonexistent")` クエリが実行される
- **THEN** Repositoryは `nil` と `error` を返す
- **AND** Resolverは GraphQL エラーレスポンスを返す
- **AND** エラーメッセージは "message not found" を含む

### Requirement: Firestore Data Operations
システムはFirestore APIを使用してドキュメントの読み取りを実行することをMUSTとする。

#### Scenario: コレクション全体の取得
- **WHEN** `List()` メソッドが呼び出される
- **THEN** `messages` コレクション全体をクエリする
- **AND** 各ドキュメントをMessage構造体にマッピングする
- **AND** `createdAt` フィールドで降順ソートする

#### Scenario: 単一ドキュメントの取得
- **WHEN** `GetByID("123")` メソッドが呼び出される
- **THEN** `messages/123` ドキュメントパスで取得する
- **AND** ドキュメントが存在すればMessage構造体にマッピングする
- **AND** 存在しなければエラーを返す

#### Scenario: Firestoreデータ型のマッピング
- **WHEN** FirestoreドキュメントをGoの構造体にマッピングする
- **THEN** Firestore `string` フィールドはGo `string` にマッピングする
- **AND** Firestore `timestamp` フィールドはGo `time.Time` を経由してISO8601文字列にマッピングする
- **AND** フィールド名はFirestoreのフィールド名と一致する

### Requirement: Sample Data Seeding
システムはFirestoreに初期サンプルデータを投入する手段を提供することをMUSTとする。

#### Scenario: シードスクリプトの存在
- **WHEN** プロジェクトディレクトリを確認する
- **THEN** `scripts/seed-firestore.go` ファイルが存在する
- **AND** スクリプトは `FIRESTORE_EMULATOR_HOST` 環境変数を読み込む

#### Scenario: サンプルデータの投入
- **WHEN** `go run scripts/seed-firestore.go` を実行する
- **THEN** 少なくとも3つのMessageドキュメントが `messages` コレクションに作成される
- **AND** 各ドキュメントは `id`, `content`, `author`, `createdAt` フィールドを持つ
- **AND** 成功メッセージがコンソールに表示される

#### Scenario: べき等性
- **WHEN** シードスクリプトを複数回実行する
- **THEN** 既存のデータを上書きまたはスキップする
- **AND** エラーが発生しない

### Requirement: Error Handling and Logging
システムはFirestore操作のエラーを適切にハンドリングし、ログに記録することをMUSTとする。

#### Scenario: 接続エラーのハンドリング
- **WHEN** Firestore Emulatorが起動していない状態でクエリを実行する
- **THEN** Repositoryは接続エラーを返す
- **AND** Resolverはエラーをログに記録する
- **AND** GraphQLエラーレスポンスをクライアントに返す

#### Scenario: 操作ログの記録
- **WHEN** Firestore操作（List, GetByID）が実行される
- **THEN** 操作の開始と完了がログに記録される
- **AND** エラーが発生した場合、エラー詳細がログに記録される

### Requirement: GraphQL Integration Testing
システムはGraphQL Playgroundを通じてFirestoreクエリをテストできることをMUSTとする。

#### Scenario: messagesクエリの実行成功
- **WHEN** GraphQL Playgroundで `{ messages { id content author createdAt } }` を実行する
- **THEN** Firestoreに存在する全Messageが返される
- **AND** レスポンスは有効なJSONである
- **AND** HTTPステータスコードは200である

#### Scenario: messageクエリの実行成功（ID指定）
- **WHEN** GraphQL Playgroundで `{ message(id: "msg1") { id content author createdAt } }` を実行する
- **THEN** ID "msg1" のMessageが返される
- **AND** レスポンスはMessageの全フィールドを含む

#### Scenario: スキーマエクスプローラーでの確認
- **WHEN** GraphQL Playgroundのスキーマエクスプローラーを開く
- **THEN** `Query.messages` が表示される
- **AND** `Query.message(id: ID!)` が表示される
- **AND** `Message` タイプとそのフィールドが表示される
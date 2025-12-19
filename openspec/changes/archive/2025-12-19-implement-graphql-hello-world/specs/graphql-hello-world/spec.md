# GraphQL Hello World

## ADDED Requirements

### Requirement: GraphQL Schema Definition
システムはGraphQLスキーマを定義し、単一のQueryタイプを含むことをMUSTとする。

#### Scenario: スキーマにhelloクエリが定義されている
- **WHEN** GraphQLスキーマファイルを読み込む
- **THEN** `Query` タイプに `hello` フィールドが存在する
- **AND** `hello` フィールドは引数を取らず `String!` を返す

### Requirement: Hello World Query
システムは `hello` クエリに対して "Hello World" を返すことをMUSTとする。

#### Scenario: helloクエリの実行成功
- **WHEN** クライアントが `{ hello }` クエリをGraphQLサーバーに送信する
- **THEN** サーバーは `"Hello World"` を含むJSONレスポンスを返す
- **AND** HTTPステータスコードは200である

#### Scenario: GraphQL Playgroundでの動作確認
- **WHEN** ブラウザでGraphQL Playgroundにアクセスする
- **THEN** スキーマエクスプローラーに `hello` クエリが表示される
- **AND** `{ hello }` クエリを実行すると `"Hello World"` が返される

### Requirement: GraphQL Server Initialization
システムはgqlgenを使用してGraphQLサーバーを起動し、HTTPリクエストを処理することをMUSTとする。

#### Scenario: サーバー起動成功
- **WHEN** `go run server.go` コマンドを実行する
- **THEN** GraphQLサーバーが指定されたポート（例: 8080）でリッスンを開始する
- **AND** コンソールに起動成功メッセージが表示される

#### Scenario: HTTPエンドポイントの提供
- **WHEN** サーバーが起動している状態で `/query` パスにアクセスする
- **THEN** GraphQL APIエンドポイントがリクエストを受け付ける
- **AND** `/` パスにアクセスするとGraphQL Playgroundが表示される

### Requirement: Code Generation with gqlgen
システムはgqlgenを使用してGraphQLスキーマからGoコードを自動生成することをMUSTとする。

#### Scenario: スキーマからコード生成
- **WHEN** `go run github.com/99designs/gqlgen generate` コマンドを実行する
- **THEN** `graph/generated.go` にGraphQLの型定義とリゾルバーインターフェースが生成される
- **AND** `graph/model/` に対応するGoの型が生成される

#### Scenario: gqlgen設定ファイルの存在
- **WHEN** プロジェクトルートを確認する
- **THEN** `gqlgen.yml` ファイルが存在する
- **AND** スキーマファイルのパス、生成先ディレクトリが正しく設定されている
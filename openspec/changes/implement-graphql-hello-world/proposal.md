# Change: GraphQL Hello World実装

## Why
GraphQL APIの基礎を学習するため、最もシンプルな「Hello World」を返すGraphQLサーバーを実装する。これがプロジェクトの土台となり、より複雑な機能を追加する前にgqlgenのセットアップとGraphQLの基本パターンを確立する。

## What Changes
- gqlgenの初期化とセットアップ
- GraphQLスキーマ定義（単一のQueryフィールド `hello`）
- GraphQLサーバーの起動と基本的なHTTPハンドラー
- `hello` クエリに対するリゾルバー実装
- gqlgen設定ファイルと生成されたコード
- 基本的な動作確認用のREADME

## Impact
- Affected specs: `graphql-hello-world` (新規追加)
- Affected code:
  - 新規: `server.go` - GraphQLサーバーのエントリーポイント
  - 新規: `graph/schema.graphqls` - GraphQLスキーマ定義
  - 新規: `graph/resolver.go` - リゾルバーの実装
  - 新規: `gqlgen.yml` - gqlgen設定ファイル
  - 新規: `go.mod` / `go.sum` - Go依存関係
# Implementation Tasks

## 1. プロジェクトセットアップ
- [x] 1.1 Goモジュールの初期化 (`go.mod` 作成)
- [x] 1.2 gqlgen依存関係のインストール
- [x] 1.3 gqlgenの初期化 (`gqlgen.yml` 生成)

## 2. GraphQLスキーマ定義
- [x] 2.1 `graph/schema.graphqls` ファイルの作成
- [x] 2.2 `Query` タイプの定義
- [x] 2.3 `hello` フィールドの追加（引数なし、`String!` 型を返す）

## 3. コード生成とリゾルバー実装
- [x] 3.1 gqlgenでコード生成を実行
- [x] 3.2 生成された `graph/resolver.go` を確認
- [x] 3.3 `hello` クエリのリゾルバー実装（"Hello World" を返す）

## 4. GraphQLサーバーの実装
- [x] 4.1 `server.go` の作成
- [x] 4.2 GraphQLハンドラーの初期化
- [x] 4.3 `/query` エンドポイントの設定
- [x] 4.4 `/` でGraphQL Playgroundを提供
- [x] 4.5 HTTPサーバーの起動処理の実装

## 5. 動作確認とドキュメント
- [x] 5.1 サーバーの起動確認
- [x] 5.2 `{ hello }` クエリの実行確認
- [x] 5.3 GraphQL Playgroundでの動作確認
- [x] 5.4 README.mdに起動手順とクエリ例を記載
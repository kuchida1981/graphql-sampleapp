# Implementation Tasks

## 1. プロジェクトセットアップ
- [ ] 1.1 Goモジュールの初期化 (`go.mod` 作成)
- [ ] 1.2 gqlgen依存関係のインストール
- [ ] 1.3 gqlgenの初期化 (`gqlgen.yml` 生成)

## 2. GraphQLスキーマ定義
- [ ] 2.1 `graph/schema.graphqls` ファイルの作成
- [ ] 2.2 `Query` タイプの定義
- [ ] 2.3 `hello` フィールドの追加（引数なし、`String!` 型を返す）

## 3. コード生成とリゾルバー実装
- [ ] 3.1 gqlgenでコード生成を実行
- [ ] 3.2 生成された `graph/resolver.go` を確認
- [ ] 3.3 `hello` クエリのリゾルバー実装（"Hello World" を返す）

## 4. GraphQLサーバーの実装
- [ ] 4.1 `server.go` の作成
- [ ] 4.2 GraphQLハンドラーの初期化
- [ ] 4.3 `/query` エンドポイントの設定
- [ ] 4.4 `/` でGraphQL Playgroundを提供
- [ ] 4.5 HTTPサーバーの起動処理の実装

## 5. 動作確認とドキュメント
- [ ] 5.1 サーバーの起動確認
- [ ] 5.2 `{ hello }` クエリの実行確認
- [ ] 5.3 GraphQL Playgroundでの動作確認
- [ ] 5.4 README.mdに起動手順とクエリ例を記載
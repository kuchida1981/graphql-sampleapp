# GraphQL Sample App

GraphQLのベストプラクティスとパターンを学習するためのチュートリアルプロジェクト。gqlgenを使用したプロダクション対応GraphQL APIの構築方法を示すリファレンス実装。

## 技術スタック

- **言語:** Go (Golang)
- **GraphQLフレームワーク:** gqlgen

## はじめに

### 前提条件

- Go 1.24.0以上

### インストール

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

### cURLでのクエリ実行

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ hello }"}'
```

## プロジェクト構造

```
.
├── server.go              # GraphQLサーバーのエントリーポイント
├── gqlgen.yml             # gqlgen設定ファイル
├── graph/
│   ├── schema.graphqls    # GraphQLスキーマ定義
│   ├── resolver.go        # リゾルバーのベース構造
│   ├── schema.resolvers.go # リゾルバー実装
│   ├── generated.go       # gqlgenが生成したコード
│   └── model/             # GraphQLモデルの型定義
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

# Project Context

## Purpose
GraphQLのベストプラクティスとパターンを学習するためのチュートリアルプロジェクト。gqlgenを使用したプロダクション対応GraphQL APIの構築方法を示し、リアルタイムサブスクリプションとマルチデータベース統合を実装したリファレンス実装。

## Tech Stack
- **言語:** Go (Golang)
- **GraphQLフレームワーク:** gqlgen - GraphQLサーバー構築用のGoライブラリ
- **データベース:**
  - PostgreSQL - メインのリレーショナルデータベース
  - Cloud Firestore - 特定のユースケース向けドキュメントストア
- **リアルタイム:** WebSocketベースのGraphQLサブスクリプション
- **開発ツール:**
  - gofmt/goimports - コード整形
  - Goの標準ツールチェーン

## Project Conventions

### Code Style
- **フォーマット:** すべてのGoコードにgofmtとgoimportsを使用
- **命名規則:** Goの命名規則に従う（エクスポートはCamelCase、非エクスポートはcamelCase）
- **パッケージ構成:** 技術レイヤーではなく、機能/ドメインごとに整理
- **エラーハンドリング:** 標準的なGoのエラーハンドリングパターンを使用
- **インポート:** インポートを3つのセクションに分ける（標準ライブラリ、外部、内部）

### Architecture Patterns
- **GraphQLスキーマファースト:** GraphQL SDLでスキーマを定義し、gqlgenでコード生成
- **Resolverパターン:** Resolverは薄く保ち、ビジネスロジックはサービス層に委譲
- **Repositoryパターン:** リポジトリインターフェースの背後にデータベースアクセスを抽象化
- **依存性注入:** コンストラクタを通じて依存関係を明示的に渡す
- **Contextの伝播:** リクエストスコープの値とキャンセル処理にcontext.Contextを使用

### Testing Strategy
- **ユニットテスト:** 分離されたコンポーネントテストにGoの標準testingパッケージを使用
- **統合テスト:** 実際のデータベース接続でGraphQL resolverをテスト
- **モッキング:** 外部依存関係にはインターフェースとモックツール（gomock）を使用
- **テストの配置:** テスト対象コードと同じパッケージにテストを配置（*_test.go）
- **カバレッジ:** パーセンテージ目標ではなく、意味のあるカバレッジを目指す

### Git Workflow
- **戦略:** トランクベース開発
- **メインブランチ:** mainブランチは常にデプロイ可能な状態
- **フィーチャーブランチ:** 短命なブランチ（2日未満）を頻繁にマージ
- **コミットメッセージ:** Conventional Commits形式（feat:, fix:, docs:など）
- **プルリクエスト:** すべての変更にコードレビュー付きPRが必須
- **CI/CD:** マージ前にすべてのブランチでテストとlintingを実行

## Domain Context
GraphQLの概念を学習するためのサンプルアプリケーション。ドメインの複雑さではなくGraphQLパターンに焦点を当てるため、ドメインモデルはシンプルで親しみやすいもの（例: ユーザー、投稿、コメント）にする。

実装すべき主なGraphQLの概念:
- Query、Mutation、Subscription
- N+1クエリ問題回避のためのDataLoaderパターン
- GraphQLコンテキストでの認証と認可
- エラーハンドリングとバリデーション
- ページネーションパターン（カーソルベース、オフセットベース）
- サブスクリプションによるリアルタイム更新

## Important Constraints
- **シンプルさ優先:** 実装はシンプルで教育的に保ち、過度なエンジニアリングを避ける
- **ドキュメント:** コードは自己文書化され、明らかでないロジックにはコメントを付ける
- **Goのイディオム:** 他言語のパターンではなく、確立されたGoのパターンとイディオムに従う
- **スキーマ設計:** GraphQLスキーマは直感的でベストプラクティスに従う
- **後方互換性:** GraphQLスキーマの変更は可能な限り後方互換性を維持する

## External Dependencies
- **PostgreSQL:** 構造化データ用のリレーショナルデータベース
- **Cloud Firestore:** ドキュメントベースデータ用のGoogle Cloud Firestore
- **gqlgen:** GraphQLサーバーライブラリとコードジェネレーター
- **データベースドライバ:**
  - pgx/v5 - PostgreSQLドライバとツールキット
  - firebase-go-sdk - Cloud Firestoreクライアント
- **認証:** 未定（JWT、OAuth2など）
- **デプロイ:** 未定（Cloud Run、Kubernetesなど）

# Proposal: PostgreSQL User Integration

## Context
現在、このGraphQLサンプルアプリケーションはFirestoreをバックエンドとしてMessage機能を実装している。プロジェクトの目的は「GraphQLのベストプラクティスとマルチデータベース統合を学ぶ」ことであり、PostgreSQLは`project.md`で既にTech Stackとして宣言されているが、まだ実装されていない。

このプロジェクトは、他サービスの改善を視野に入れたチュートリアルプロジェクトであり、そのサービスはデータ返却のみを行う。したがって、Read操作のみを実装する。

## Why
このチュートリアルプロジェクトの主要な学習目標の一つは「マルチデータベース統合」である。現在Firestoreのみが実装されているため、PostgreSQLを追加することで:

1. **異なるデータベースパラダイムの比較学習**: NoSQL（Firestore）とRDBMS（PostgreSQL）の使い分けパターンを実践的に学べる
2. **リポジトリパターンの有効性実証**: 同じインターフェースで異なるデータソースを抽象化する設計パターンを体験できる
3. **実サービスへの応用準備**: 実際のサービスでは正規化・関係性が重要なデータをPostgreSQLに、外部スキーマの影響を受けやすいデータをFirestoreに配置する設計を将来実装する基盤となる
4. **GraphQL統合の柔軟性**: 単一のGraphQLスキーマで複数のバックエンドを透過的に扱うアーキテクチャを実現する

特に、実際のサービス改善を見据えたチュートリアルとして、データ返却（Read操作）のみを実装することで、複雑さを抑えつつ本質的な統合パターンを学習できる。

## Objective
PostgreSQLをバックエンドとして、Userドメインモデルの読み取り機能をGraphQL API経由で提供する。Firestoreの実装パターン（リポジトリパターン、依存性注入、GraphQLスキーマ拡張）を踏襲し、複数データベースの統合パターンを学習可能にする。

## Scope

### In Scope
- PostgreSQLコンテナをdocker-composeに追加
- Userドメインモデルの定義（ID、名前、メール、作成日時）
- `database/sql` + `pgx/v5`ドライバを使用したPostgreSQL接続
- リポジトリパターンによるデータアクセス層の実装（Read操作のみ: List, GetByID）
- GraphQLスキーマへのUser型とクエリ追加（`users`, `user(id: ID!)`）
- 初期化SQLスクリプトによるスキーマ定義
- サンプルデータ投入用のシードスクリプト
- GraphQL Playgroundでの動作確認

### Out of Scope
- Create/Update/Delete操作（GraphQL Mutation）
- ユーザー認証・認可機能
- UserとMessageの関連づけ（将来の拡張として想定されているが、今回は扱わない）
- データベースマイグレーションツール（シンプルな初期化スクリプトのみ）
- トランザクション処理
- ページネーション（将来的な学習項目だが、今回はシンプルな全件取得のみ）

## Assumptions
- Docker Composeが利用可能
- Firestoreの実装パターン（`internal/repository`、`internal/domain`、依存性注入）を継承
- PostgreSQLは開発環境ではDockerコンテナとして起動
- エラーハンドリングとロギングはFirestore実装と同等レベル
- GraphQLスキーマの後方互換性を維持（既存のMessage関連機能に影響を与えない）

## Risks and Mitigation
- **リスク**: PostgreSQL接続が失敗した場合、アプリケーション全体が起動しない
  - **軽減策**: 接続エラー時にはログを出力し、gracefulに終了する（Firestore実装と同じパターン）
- **リスク**: 複数のリポジトリ（Firestore, PostgreSQL）を持つことで、Resolverの依存性注入が複雑になる
  - **軽減策**: Resolverに各リポジトリを個別フィールドとして持たせ、NewResolverで明示的に注入する
- **リスク**: 初期化SQLスクリプトの実行タイミングが不明確
  - **軽減策**: READMEに手動実行手順を記載し、Docker Compose起動時の初期化スクリプトとして配置

## Success Metrics
- GraphQL Playgroundで`{ users { id name email createdAt } }`クエリが正常に実行できる
- GraphQL Playgroundで`{ user(id: "1") { id name email createdAt } }`クエリが正常に実行できる
- `docker compose up`でPostgreSQLとアプリケーションが正常に起動する
- シードスクリプト実行後、サンプルユーザーデータが取得できる
- 既存のFirestore Message機能に影響がない（`messages`クエリが引き続き動作）

## Related Changes
なし（初回のPostgreSQL統合）

## Notes
- 将来的には、UserとMessageを組み合わせた機能（例: 特定ユーザーのメッセージ一覧）を実装する可能性がある
- ページネーション、フィルタリング、ソートなどの高度なクエリ機能は将来の学習項目として残す
- この変更はFirestore実装と同じパターンを意図的に踏襲し、複数データベース統合の学習を容易にする
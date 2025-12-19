# Proposal: Add Firestore Docker Integration

## Overview
このプロポーザルは、既存のGraphQL Hello Worldサーバーに対して、Cloud Firestoreをバックエンドとして追加し、Docker Composeによるローカル開発環境を整備します。Firestoreからデータを取得し返却するサンプル機能を実装することで、マルチデータベース統合とドキュメントストアの使用パターンを学習できる環境を提供します。

## Why
現在のプロジェクトはGraphQL Hello Worldのみを実装しており、実際のデータストアとの連携がありません。このままでは、GraphQLとデータベースの統合パターンやリポジトリパターンの実装方法を学習できません。

Cloud Firestoreをバックエンドとして追加し、Docker Composeによるローカル開発環境を整備することで、以下の学習目標を達成します：

- **マルチデータベースアーキテクチャの学習**: PostgreSQL（将来追加予定）とFirestoreを併用するパターンの基礎を構築
- **ドキュメントストアの活用**: NoSQLデータベースとGraphQLの統合方法を実践的に学習
- **ローカル開発環境の整備**: Docker Composeで簡単に起動できる開発環境により、セットアップの障壁を下げる
- **リアルワールドなサンプル**: 単なるHello Worldから一歩進んだ、実践的なデータ取得実装を通じてベストプラクティスを習得

また、Docker Compose環境を整備することで、チーム開発や継続的インテグレーションの基盤も同時に構築できます。

## Scope
このプロポーザルは以下を含みます：

### In Scope
1. **Docker Compose設定**
   - Goアプリケーション（GraphQLサーバー）用のコンテナ
   - Cloud Firestore Emulator用のコンテナ
   - コンテナ間ネットワーキング設定

2. **Firestore統合**
   - Firebase Admin SDKのセットアップ
   - Firestoreクライアント初期化コード
   - リポジトリパターンによるFirestoreアクセス層

3. **サンプル機能実装**
   - Firestoreからデータを取得するGraphQLクエリ
   - サンプルデータのシードスクリプト
   - リゾルバー実装

4. **開発者体験向上**
   - READMEへのDocker Compose使用方法追記
   - サンプルデータのセットアップ手順

### Out of Scope
- PostgreSQLの統合（別のプロポーザルで対応）
- 認証・認可機能
- Mutation操作（将来のプロポーザルで対応）
- プロダクション環境へのデプロイ設定
- Firestore セキュリティルール
- 複雑なデータモデリング（シンプルなサンプルのみ）

## Deliverables
1. `docker-compose.yml` - アプリケーションとFirestore Emulatorの定義
2. `Dockerfile` - Goアプリケーション用のコンテナイメージ定義
3. Firestore統合コード（リポジトリ、クライアント初期化）
4. 更新されたGraphQLスキーマ（Firestoreからデータを取得するクエリを含む）
5. サンプルデータシードスクリプト
6. 更新されたREADME（Docker Composeセットアップ手順を含む）

## Success Criteria
- `docker-compose up` コマンドでアプリケーションとFirestore Emulatorが起動する
- GraphQL PlaygroundでFirestoreからデータを取得するクエリが実行できる
- サンプルデータが正しく返される
- コードがGoの規約とプロジェクトのアーキテクチャパターンに準拠している
- READMEの手順に従って開発環境をセットアップできる

## Dependencies
- 既存の `graphql-hello-world` spec（基盤となるGraphQLサーバー）
- Docker および Docker Compose（開発環境要件）
- Firebase Admin SDK for Go

## Assumptions
- 開発者のローカル環境にDocker/Docker Composeがインストール済み
- Firestore Emulatorを使用し、実際のGCPプロジェクトは不要
- サンプルデータは簡素なドメインモデル（例: メッセージ、ユーザーなど）
- 環境変数でFirestoreエンドポイント（Emulator vs 本番）を切り替え可能にする
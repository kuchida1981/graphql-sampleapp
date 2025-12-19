# Proposal: Add Testing Framework

## Overview
現在のコードベースにはテストコードが一切存在しない状態です。本提案では、既存の実装に対する包括的なテストスイートを追加し、プロジェクト全体のテストポリシーを確立します。テストはモックを使用してDBへの依存を排除し、CI/CDパイプラインで自動実行可能な形で実装します。

## Motivation
- **品質保証**: リファクタリングや機能追加時のリグレッション検出
- **ドキュメント**: テストコードが実装の使用例とドキュメントの役割を果たす
- **開発速度**: テストがあることで安心してリファクタリング・変更ができる
- **教育的価値**: チュートリアルプロジェクトとして、Goのテストベストプラクティスを示す

## Scope
このプロポーザルでは以下を対象とします:

### 含まれるもの
1. **Repositoryレイヤーのテスト**
   - `PostgresUserRepository`
   - `FirestoreMessageRepository`
   - `PostgresWeatherAlertMetadataRepository`
   - `FirestoreWeatherAlertRepository`

2. **GraphQL Resolverのテスト**
   - `Hello` resolver
   - `Users` / `User` resolvers
   - `Messages` / `Message` resolvers
   - `WeatherAlerts` resolver

3. **テストポリシードキュメント**
   - テスト作成の基本方針
   - モッキング戦略
   - カバレッジ目標
   - CI/CD統合

### 含まれないもの
- 統合テスト (実際のDBを使用するテスト)
- E2Eテスト
- パフォーマンステスト
- セキュリティテスト

## Dependencies
- 既存のすべてのspec (既存実装に対するテストを追加)
- Goの標準`testing`パッケージ
- モッキングライブラリ (sqlmock, gomock等)

## Breaking Changes
なし。既存のコードには変更を加えず、テストコードのみを追加します。

## Alternatives Considered
1. **統合テストから始める**: 実際のDBを使う統合テストを先に実装する案
   - 却下理由: セットアップが複雑で実行が遅い。まずは高速なユニットテストを確立すべき

2. **テストフレームワーク (testify等) を使用**: サードパーティのテストフレームワークを導入する案
   - 却下理由: プロジェクトポリシーでシンプルさを優先。標準ライブラリで十分

## Success Criteria
- [ ] すべてのRepository実装に対するユニットテストが存在する
- [ ] すべてのGraphQL Resolverに対するユニットテストが存在する
- [ ] テストポリシードキュメントが作成され、プロジェクトルートに配置される
- [ ] `go test ./...` がすべてパスする
- [ ] モックを使用してDBへの依存がない
- [ ] CI/CDでテストが自動実行される設定がある
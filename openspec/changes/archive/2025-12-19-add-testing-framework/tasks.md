# Tasks: Add Testing Framework

## Task 1: Set up testing dependencies
- [x] `go get github.com/DATA-DOG/go-sqlmock`をインストール
- [x] `go get github.com/golang/mock/gomock`をインストール (オプション)
- [x] `go mod tidy`で依存関係を整理
- **Validation**: `go.mod`に両パッケージが追加されている
- **Deliverable**: ユーザーは`go test`を実行できる環境が整う

## Task 2: Create PostgreSQL User Repository tests
- [x] `internal/repository/postgres_user_test.go`を作成
- [x] `TestPostgresUserRepository_List`を実装
  - 正常系: ユーザーリスト取得成功
  - 異常系: クエリエラー
  - 異常系: スキャンエラー
- [x] `TestPostgresUserRepository_GetByID`を実装
  - 正常系: ユーザー取得成功
  - 異常系: ユーザーが見つからない (sql.ErrNoRows)
  - 異常系: スキャンエラー
- **Validation**: `go test ./internal/repository -run TestPostgresUserRepository`がパスする
- **Deliverable**: PostgreSQL User Repositoryの全メソッドにテストが存在

## Task 3: Create Firestore Message Repository tests
- [x] `internal/repository/firestore_message_test.go`を作成
- [x] Firestoreクライアントのモック実装を作成
- [x] `TestFirestoreMessageRepository_List`を実装
  - 正常系: メッセージリスト取得成功
  - 異常系: イテレーションエラー
  - 異常系: データ変換エラー
- [x] `TestFirestoreMessageRepository_GetByID`を実装
  - 正常系: メッセージ取得成功
  - 異常系: ドキュメントが見つからない
  - 異常系: データ変換エラー
- **Validation**: `go test ./internal/repository -run TestFirestoreMessageRepository`がパスする
- **Deliverable**: Firestore Message Repositoryの全メソッドにテストが存在

## Task 4: Create PostgreSQL Weather Alert Metadata Repository tests
- [x] `internal/repository/postgres_weather_alert_metadata_test.go`を作成
- [x] `TestPostgresWeatherAlertMetadataRepository_Search`を実装
  - 正常系: フィルタなしで検索
  - 正常系: 地域フィルタ付き検索
  - 正常系: 日時フィルタ付き検索
  - 正常系: 複合フィルタ検索
  - 異常系: クエリエラー
- [x] `TestPostgresWeatherAlertMetadataRepository_SearchIDs`を実装
  - 正常系: IDリスト取得
  - 異常系: クエリエラー
- **Validation**: `go test ./internal/repository -run TestPostgresWeatherAlertMetadataRepository`がパスする
- **Deliverable**: Weather Alert Metadata Repositoryの全メソッドにテストが存在

## Task 5: Create Firestore Weather Alert Repository tests
- [x] `internal/repository/firestore_weather_alert_test.go`を作成
- [x] `TestFirestoreWeatherAlertRepository_GetByIDs`を実装
  - 正常系: 複数ID指定で取得
  - 正常系: 空のIDリスト
  - 正常系: 一部のIDが存在しないケース
  - 異常系: Firestoreエラー
- **Validation**: `go test ./internal/repository -run TestFirestoreWeatherAlertRepository`がパスする
- **Deliverable**: Weather Alert Repositoryの全メソッドにテストが存在

## Task 6: Create GraphQL Resolver tests - Basic queries
- [x] `graph/schema.resolvers_test.go`を作成
- [x] Repositoryインターフェースのモック実装を作成
  - `mockUserRepository`
  - `mockMessageRepository`
  - `mockWeatherAlertRepository`
  - `mockWeatherAlertMetadataRepository`
- [x] `TestQueryResolver_Hello`を実装
- [x] `TestQueryResolver_Users`を実装
  - 正常系: ユーザーリスト取得
  - 異常系: Repositoryエラー
- [x] `TestQueryResolver_User`を実装
  - 正常系: 単一ユーザー取得
  - 異常系: ユーザーが見つからない
- [x] `TestQueryResolver_Messages`を実装
  - 正常系: メッセージリスト取得
  - 異常系: Repositoryエラー
- [x] `TestQueryResolver_Message`を実装
  - 正常系: 単一メッセージ取得
  - 異常系: メッセージが見つからない
- **Validation**: `go test ./graph -run 'TestQueryResolver_(Hello|Users|User|Messages|Message)'`がパスする
- **Deliverable**: 基本的なGraphQL queriesのテストが存在

## Task 7: Create GraphQL Resolver tests - Weather alerts
- [x] `TestQueryResolver_WeatherAlerts`を実装
  - 正常系: フィルタなしで取得
  - 正常系: 地域フィルタ付き取得
  - 正常系: 日時フィルタ付き取得
  - 正常系: メタデータが見つからない (空リスト返却)
  - 正常系: Firestoreデータが一部欠損
  - 異常系: 不正な日時フォーマット
  - 異常系: Metadataリポジトリエラー
  - 異常系: Firestoreリポジトリエラー
- **Validation**: `go test ./graph -run TestQueryResolver_WeatherAlerts`がパスする
- **Deliverable**: WeatherAlerts queryの全シナリオにテストが存在

## Task 8: Create TESTING.md documentation
- [x] プロジェクトルートに`TESTING.md`を作成
- [x] 以下のセクションを記述:
  - テスト作成の基本方針
  - モッキング戦略
  - テスト命名規則
  - テーブル駆動テストの書き方
  - カバレッジ方針
  - テスト実行方法
  - CI/CD統合
  - 実際のテストコード例
  - 必要な依存ライブラリ
- [x] すべての説明を日本語で記述
- **Validation**: `TESTING.md`が存在し、すべてのセクションが完成している
- **Deliverable**: チーム全員がテスト作成方法を理解できるドキュメント

## Task 9: Run all tests and verify
- [x] `go test ./...`を実行してすべてのテストがパスすることを確認
- [x] `go test -cover ./...`でカバレッジを確認
- [x] カバレッジレポートを確認し、主要なコードパスがカバーされていることを検証
- **Validation**: すべてのテストが成功し、エラーがない
- **Deliverable**: 完全に動作するテストスイート

## Task 10: Add CI/CD configuration (optional, recommended)
- [x] `.github/workflows/test.yml`を作成 (GitHub Actions使用の場合)
- [x] 以下のステップを定義:
  - Go環境のセットアップ
  - 依存関係のインストール (`go mod download`)
  - テスト実行 (`go test ./...`)
  - カバレッジレポート生成 (オプション)
- [x] PRで自動実行されることを確認
- **Validation**: PRを作成してCI/CDでテストが実行される
- **Deliverable**: 自動化されたテスト実行環境

---

## Dependencies Between Tasks
- Task 1は他のすべてのタスクの前提条件
- Task 2-5は並行実行可能 (Repository tests)
- Task 6-7はTask 2-5の後が望ましい (Repositoryの動作を確認してから)
- Task 8は独立して実行可能
- Task 9はTask 2-7の完了後
- Task 10は独立して実行可能だが、Task 9の前が望ましい

## Estimated Effort
- Task 1: 5分 (依存関係インストール)
- Task 2: 1-2時間 (PostgreSQL User Repository tests)
- Task 3: 1-2時間 (Firestore Message Repository tests)
- Task 4: 1.5-2時間 (PostgreSQL Weather Alert Metadata tests)
- Task 5: 1-1.5時間 (Firestore Weather Alert tests)
- Task 6: 2-3時間 (Basic resolver tests)
- Task 7: 1.5-2時間 (Weather alerts resolver tests)
- Task 8: 1-1.5時間 (Documentation)
- Task 9: 30分 (Verification)
- Task 10: 30分-1時間 (CI/CD setup)

**Total**: 約10-15時間

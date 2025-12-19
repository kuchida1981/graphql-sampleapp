# Spec: Repository Unit Tests

## ADDED Requirements

### Requirement: PostgreSQL User Repository Tests
すべての`PostgresUserRepository`メソッドに対して、モックDBを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: List users successfully
- **Given**: モックDBが有効なユーザーデータを返す設定
- **When**: `List(ctx)`を呼び出す
- **Then**: 期待されるユーザーリストが返される

#### Scenario: List users with database error
- **Given**: モックDBがエラーを返す設定
- **When**: `List(ctx)`を呼び出す
- **Then**: エラーが返される

#### Scenario: Get user by ID successfully
- **Given**: モックDBが指定IDのユーザーデータを返す設定
- **When**: `GetByID(ctx, "user-1")`を呼び出す
- **Then**: 該当ユーザーが返される

#### Scenario: Get user by ID when not found
- **Given**: モックDBが`sql.ErrNoRows`を返す設定
- **When**: `GetByID(ctx, "non-existent")`を呼び出す
- **Then**: "user not found"エラーが返される

---

### Requirement: Firestore Message Repository Tests
すべての`FirestoreMessageRepository`メソッドに対して、モックFirestoreクライアントを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: List messages successfully
- **Given**: モックFirestoreが有効なメッセージデータを返す設定
- **When**: `List(ctx)`を呼び出す
- **Then**: 期待されるメッセージリストが返される

#### Scenario: List messages with Firestore error
- **Given**: モックFirestoreがエラーを返す設定
- **When**: `List(ctx)`を呼び出す
- **Then**: エラーが返される

#### Scenario: Get message by ID successfully
- **Given**: モックFirestoreが指定IDのメッセージデータを返す設定
- **When**: `GetByID(ctx, "msg-1")`を呼び出す
- **Then**: 該当メッセージが返される

#### Scenario: Get message by ID when not found
- **Given**: モックFirestoreがドキュメント未存在エラーを返す設定
- **When**: `GetByID(ctx, "non-existent")`を呼び出す
- **Then**: "message not found"エラーが返される

---

### Requirement: PostgreSQL Weather Alert Metadata Repository Tests
すべての`PostgresWeatherAlertMetadataRepository`メソッドに対して、モックDBを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: Search metadata without filters
- **Given**: モックDBがすべてのメタデータを返す設定
- **When**: `Search(ctx, MetadataFilter{})`を呼び出す
- **Then**: すべてのメタデータが発行日時降順で返される

#### Scenario: Search metadata with region filter
- **Given**: モックDBが特定地域のメタデータを返す設定
- **When**: `Search(ctx, MetadataFilter{Region: &"Tokyo"})`を呼び出す
- **Then**: 東京地域のメタデータのみが返される

#### Scenario: Search metadata with issued after filter
- **Given**: モックDBが指定日時以降のメタデータを返す設定
- **When**: `Search(ctx, MetadataFilter{IssuedAfter: &time})`を呼び出す
- **Then**: 指定日時以降のメタデータのみが返される

#### Scenario: Search metadata with combined filters
- **Given**: モックDBが地域と日時の両方でフィルタされたメタデータを返す設定
- **When**: `Search(ctx, MetadataFilter{Region: &"Osaka", IssuedAfter: &time})`を呼び出す
- **Then**: 両条件を満たすメタデータのみが返される

#### Scenario: SearchIDs returns only IDs
- **Given**: モックDBがIDリストを返す設定
- **When**: `SearchIDs(ctx, filter)`を呼び出す
- **Then**: IDの文字列スライスのみが返される

---

### Requirement: Firestore Weather Alert Repository Tests
すべての`FirestoreWeatherAlertRepository`メソッドに対して、モックFirestoreクライアントを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: Get weather alerts by multiple IDs
- **Given**: モックFirestoreが複数のアラートデータを返す設定
- **When**: `GetByIDs(ctx, []string{"alert-1", "alert-2"})`を呼び出す
- **Then**: 指定IDに対応するすべてのアラートが返される

#### Scenario: Get weather alerts with partial missing data
- **Given**: モックFirestoreで一部のIDが存在しない設定
- **When**: `GetByIDs(ctx, []string{"alert-1", "non-existent"})`を呼び出す
- **Then**: 存在するアラートのみが返される (エラーにならない)

#### Scenario: Get weather alerts with empty ID list
- **Given**: 空のIDスライス
- **When**: `GetByIDs(ctx, []string{})`を呼び出す
- **Then**: 空のスライスが返される

---

### Requirement: Test File Organization
テストファイルは対象コードと同じパッケージに配置され、`_test.go`サフィックスを持たなければならない (MUST have)。

#### Scenario: Repository test file placement
- **Given**: `internal/repository/postgres_user.go`が存在する
- **When**: テストファイルを作成する
- **Then**: `internal/repository/postgres_user_test.go`として作成される

---

### Requirement: Mock Database Dependencies
テストは実際のデータベース接続を使用せず、モックライブラリで代替しなければならない (MUST use mocks)。

#### Scenario: PostgreSQL mocking with sqlmock
- **Given**: PostgreSQL Repositoryのテスト
- **When**: テストをセットアップする
- **Then**: `github.com/DATA-DOG/go-sqlmock`を使用してDB接続をモック化する

#### Scenario: No real database connections in unit tests
- **Given**: すべてのRepositoryユニットテスト
- **When**: テストを実行する
- **Then**: 実際のPostgreSQLまたはFirestoreへの接続は発生しない

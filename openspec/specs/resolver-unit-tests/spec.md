# resolver-unit-tests Specification

## Purpose
TBD - created by archiving change add-testing-framework. Update Purpose after archive.
## Requirements
### Requirement: Hello Resolver Test
`Hello` GraphQL resolverに対するユニットテストが存在しなければならない (MUST exist)。

#### Scenario: Return hello message
- **Given**: Hello resolverが呼び出される
- **When**: `Hello(ctx)`を実行する
- **Then**: "Hello World"が返される

---

### Requirement: Users Query Resolver Tests
`Users` GraphQL resolverに対して、モックRepositoryを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: Fetch users successfully
- **Given**: モックUserRepositoryが有効なユーザーリストを返す設定
- **When**: `Users(ctx)`を呼び出す
- **Then**: GraphQL model形式のユーザーリストが返される
- **And**: domain.Userからmodel.Userへの変換が正しく行われる
- **And**: CreatedAtがISO8601形式で変換される

#### Scenario: Users query with repository error
- **Given**: モックUserRepositoryがエラーを返す設定
- **When**: `Users(ctx)`を呼び出す
- **Then**: "failed to fetch users"を含むエラーが返される

---

### Requirement: User Query Resolver Tests
`User` GraphQL resolverに対して、モックRepositoryを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: Fetch single user by ID
- **Given**: モックUserRepositoryが指定IDのユーザーを返す設定
- **When**: `User(ctx, "user-1")`を呼び出す
- **Then**: GraphQL model形式のユーザーが返される
- **And**: domain.Userからmodel.Userへの変換が正しく行われる

#### Scenario: User query with non-existent ID
- **Given**: モックUserRepositoryが"user not found"エラーを返す設定
- **When**: `User(ctx, "non-existent")`を呼び出す
- **Then**: "failed to fetch user"を含むエラーが返される

---

### Requirement: Messages Query Resolver Tests
`Messages` GraphQL resolverに対して、モックRepositoryを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: Fetch messages successfully
- **Given**: モックMessageRepositoryが有効なメッセージリストを返す設定
- **When**: `Messages(ctx)`を呼び出す
- **Then**: GraphQL model形式のメッセージリストが返される
- **And**: domain.Messageからmodel.Messageへの変換が正しく行われる

#### Scenario: Messages query with repository error
- **Given**: モックMessageRepositoryがエラーを返す設定
- **When**: `Messages(ctx)`を呼び出す
- **Then**: "failed to fetch messages"を含むエラーが返される

---

### Requirement: Message Query Resolver Tests
`Message` GraphQL resolverに対して、モックRepositoryを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: Fetch single message by ID
- **Given**: モックMessageRepositoryが指定IDのメッセージを返す設定
- **When**: `Message(ctx, "msg-1")`を呼び出す
- **Then**: GraphQL model形式のメッセージが返される

#### Scenario: Message query with non-existent ID
- **Given**: モックMessageRepositoryが"message not found"エラーを返す設定
- **When**: `Message(ctx, "non-existent")`を呼び出す
- **Then**: "failed to fetch message"を含むエラーが返される

---

### Requirement: WeatherAlerts Query Resolver Tests
`WeatherAlerts` GraphQL resolverに対して、モックRepositoryを使用したユニットテストが存在しなければならない (MUST exist)。

#### Scenario: Fetch weather alerts without filters
- **Given**: モックWeatherAlertMetadataRepositoryとWeatherAlertRepositoryが有効なデータを返す設定
- **When**: `WeatherAlerts(ctx, nil, nil)`を呼び出す
- **Then**: すべてのアラートが返される
- **And**: MetadataとFirestoreデータが正しく結合される

#### Scenario: Fetch weather alerts with region filter
- **Given**: モックリポジトリが東京地域のデータのみを返す設定
- **When**: `WeatherAlerts(ctx, &"Tokyo", nil)`を呼び出す
- **Then**: 東京地域のアラートのみが返される

#### Scenario: Fetch weather alerts with issuedAfter filter
- **Given**: モックリポジトリが指定日時以降のデータを返す設定
- **When**: `WeatherAlerts(ctx, nil, &"2024-01-01T00:00:00Z")`を呼び出す
- **Then**: 指定日時以降のアラートのみが返される

#### Scenario: WeatherAlerts with invalid issuedAfter format
- **Given**: 不正なISO8601形式の日時文字列
- **When**: `WeatherAlerts(ctx, nil, &"invalid-date")`を呼び出す
- **Then**: "invalid issuedAfter format"を含むエラーが返される

#### Scenario: WeatherAlerts with no metadata found
- **Given**: モックWeatherAlertMetadataRepositoryが空のリストを返す設定
- **When**: `WeatherAlerts(ctx, &"NonExistentRegion", nil)`を呼び出す
- **Then**: 空のスライスが返される (エラーにならない)

#### Scenario: WeatherAlerts with metadata but missing Firestore data
- **Given**: メタデータは存在するが、Firestoreから対応データが取得できない設定
- **When**: `WeatherAlerts(ctx, nil, nil)`を呼び出す
- **Then**: Firestoreデータが存在するアラートのみが返される
- **And**: ログに警告が出力される

---

### Requirement: Mock Repository Dependencies
Resolverテストは実際のRepositoryを使用せず、モックインターフェースで代替しなければならない (MUST use mocks)。

#### Scenario: Create mock user repository
- **Given**: User resolverのテスト
- **When**: テストをセットアップする
- **Then**: `UserRepository`インターフェースのモック実装を作成する

#### Scenario: Create mock message repository
- **Given**: Message resolverのテスト
- **When**: テストをセットアップする
- **Then**: `MessageRepository`インターフェースのモック実装を作成する

#### Scenario: Create mock weather alert repositories
- **Given**: WeatherAlerts resolverのテスト
- **When**: テストをセットアップする
- **Then**: `WeatherAlertRepository`と`WeatherAlertMetadataRepository`の両方のモック実装を作成する

---

### Requirement: Test File Organization
Resolverテストファイルは`graph/`ディレクトリに配置され、`_test.go`サフィックスを持たなければならない (MUST have)。

#### Scenario: Resolver test file placement
- **Given**: `graph/schema.resolvers.go`が存在する
- **When**: テストファイルを作成する
- **Then**: `graph/schema.resolvers_test.go`として作成される


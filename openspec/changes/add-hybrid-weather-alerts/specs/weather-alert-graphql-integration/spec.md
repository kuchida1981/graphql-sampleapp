# weather-alert-graphql-integration Specification Delta

## ADDED Requirements

### Requirement: WeatherAlert GraphQL Type Definition
システムは気象アラートを表すGraphQLタイプを定義することをMUSTとする。

#### Scenario: WeatherAlertタイプの定義
- **WHEN** `graph/schema.graphqls` を確認する
- **THEN** `WeatherAlert` タイプが定義されている
- **AND** 以下のフィールドが含まれる:
  - `id: ID!`
  - `region: String!`
  - `severity: String!`
  - `issuedAt: String!`
  - `title: String!`
  - `description: String!`
  - `rawData: String!`
  - `affectedAreas: [String!]!`
  - `recommendations: [String!]!`

#### Scenario: フィールドの由来の明確化
- **WHEN** GraphQLスキーマを確認する
- **THEN** `region`, `severity`, `issuedAt` はPostgreSQLメタデータから取得される
- **AND** `title`, `description`, `rawData`, `affectedAreas`, `recommendations` はFirestoreから取得される
- **AND** `id` は両方のデータストアで共通のキーである

### Requirement: weatherAlerts Query Definition
システムは気象アラートを検索するGraphQLクエリを定義することをMUSTとする。

#### Scenario: weatherAlertsクエリの定義
- **WHEN** `graph/schema.graphqls` の `Query` タイプを確認する
- **THEN** `weatherAlerts(region: String, issuedAfter: String): [WeatherAlert!]!` クエリが定義されている
- **AND** `region` 引数はオプショナル（`String`）である
- **AND** `issuedAfter` 引数はオプショナル（`String`）であり、ISO8601形式の日時文字列を期待する
- **AND** 戻り値は非nullの `WeatherAlert` 配列である

#### Scenario: 既存クエリとの共存
- **WHEN** `graph/schema.graphqls` を確認する
- **THEN** 既存の `hello`, `messages`, `message`, `users`, `user` クエリが維持されている
- **AND** `weatherAlerts` クエリが追加されている

### Requirement: Resolver Integration with Hybrid Repositories
システムはPostgreSQLメタデータとFirestoreデータを組み合わせてGraphQLレスポンスを生成することをMUSTとする。

#### Scenario: Resolverの依存性注入
- **WHEN** `graph/resolver.go` を確認する
- **THEN** `Resolver` 構造体に以下のフィールドが追加されている:
  - `weatherAlertMetadataRepo repository.WeatherAlertMetadataRepository`
  - `weatherAlertRepo repository.WeatherAlertRepository`

#### Scenario: NewResolverコンストラクタの更新
- **WHEN** `graph/resolver.go` の `NewResolver` を確認する
- **THEN** `weatherAlertMetadataRepo` と `weatherAlertRepo` を引数として受け取る
- **AND** Resolverフィールドに設定される

### Requirement: weatherAlerts Resolver Implementation
システムはハイブリッドデータフローでweatherAlertsクエリを解決することをMUSTとする。

#### Scenario: weatherAlertsリゾルバーの実装フロー
- **WHEN** `graph/schema.resolvers.go` の `WeatherAlerts` メソッドを確認する
- **THEN** 以下の処理フローが実装されている:
  1. クエリ引数から `MetadataFilter` を構築
  2. `weatherAlertMetadataRepo.SearchIDs(ctx, filter)` を呼び出し、ID配列を取得
  3. `weatherAlertRepo.GetByIDs(ctx, ids)` を呼び出し、Firestoreデータを取得
  4. PostgreSQLメタデータとFirestoreデータをマージして `WeatherAlert` GraphQLモデルに変換
  5. GraphQLレスポンスとして返す

#### Scenario: weatherAlertsリゾルバーの実装（フィルタなし）
- **WHEN** `{ weatherAlerts { id region severity } }` クエリが実行される
- **THEN** `SearchIDs(ctx, MetadataFilter{})` が呼び出される
- **AND** 全ての気象アラートIDが取得される
- **AND** Firestoreから対応するデータが取得される
- **AND** 全ての気象アラートがレスポンスに含まれる

#### Scenario: weatherAlertsリゾルバーの実装（地域フィルタあり）
- **WHEN** `{ weatherAlerts(region: "Tokyo") { id region } }` クエリが実行される
- **THEN** `SearchIDs(ctx, MetadataFilter{Region: &"Tokyo"})` が呼び出される
- **AND** 地域が "Tokyo" のIDのみが取得される
- **AND** Firestoreから該当IDのデータが取得される

#### Scenario: weatherAlertsリゾルバーの実装（日時フィルタあり）
- **WHEN** `{ weatherAlerts(issuedAfter: "2025-12-19T00:00:00Z") { id issuedAt } }` クエリが実行される
- **THEN** `issuedAfter` 引数が `time.Time` にパースされる
- **AND** `SearchIDs(ctx, MetadataFilter{IssuedAfter: &parsedTime})` が呼び出される
- **AND** 指定日時以降のIDのみが取得される

#### Scenario: weatherAlertsリゾルバーの実装（複数フィルタあり）
- **WHEN** `{ weatherAlerts(region: "Tokyo", issuedAfter: "2025-12-19T00:00:00Z") { id } }` クエリが実行される
- **THEN** `SearchIDs(ctx, MetadataFilter{Region: &"Tokyo", IssuedAfter: &parsedTime})` が呼び出される
- **AND** 両方の条件を満たすIDのみが取得される

### Requirement: Data Merging Logic
システムはPostgreSQLメタデータとFirestoreデータを正しくマージすることをMUSTとする。

#### Scenario: メタデータとFirestoreデータのマージ
- **WHEN** Resolverでデータをマージする
- **THEN** PostgreSQL から取得した `id` に対応する Firestore データを検索する
- **AND** 両方のデータを組み合わせて `model.WeatherAlert` GraphQLモデルを生成する
- **AND** `rawData` フィールドは `json.Marshal()` でJSON文字列に変換される

#### Scenario: Firestoreデータが見つからない場合のハンドリング
- **WHEN** PostgreSQLにIDが存在するがFirestoreにデータが存在しない
- **THEN** 該当IDはGraphQLレスポンスから除外される
- **AND** 警告ログが記録される
- **AND** エラーは返されない（部分的な成功を許容）

#### Scenario: 空の結果のハンドリング
- **WHEN** フィルタ条件に合致するデータが存在しない
- **THEN** 空の配列 `[]` がGraphQLレスポンスとして返される
- **AND** エラーは返されない

### Requirement: Error Handling in Resolver
システムはリゾルバーでのエラーを適切にハンドリングすることをMUSTとする。

#### Scenario: PostgreSQLエラーのハンドリング
- **WHEN** `SearchIDs` でPostgreSQLエラーが発生する
- **THEN** Resolverは GraphQL エラーレスポンスを返す
- **AND** エラーメッセージは "failed to search weather alerts" を含む
- **AND** エラーログが記録される

#### Scenario: Firestoreエラーのハンドリング
- **WHEN** `GetByIDs` でFirestoreエラーが発生する
- **THEN** Resolverは GraphQL エラーレスポンスを返す
- **AND** エラーメッセージは "failed to get weather alert details" を含む
- **AND** エラーログが記録される

#### Scenario: 日時パースエラーのハンドリング
- **WHEN** `issuedAfter` 引数が不正な形式の場合
- **THEN** Resolverは GraphQL エラーレスポンスを返す
- **AND** エラーメッセージは "invalid issuedAfter format, expected ISO8601" を含む

### Requirement: GraphQL Playground Integration Testing
システムはGraphQL Playgroundを通じて気象アラートクエリをテストできることをMUSTとする。

#### Scenario: weatherAlertsクエリの実行成功
- **WHEN** GraphQL Playgroundで `{ weatherAlerts { id region severity title } }` を実行する
- **THEN** データベースに存在する全気象アラートが返される
- **AND** レスポンスは有効なJSONである
- **AND** HTTPステータスコードは200である

#### Scenario: weatherAlertsクエリのフィルタ実行成功
- **WHEN** GraphQL Playgroundで `{ weatherAlerts(region: "Tokyo") { id region } }` を実行する
- **THEN** 地域が "Tokyo" の気象アラートのみが返される
- **AND** レスポンスに他の地域のデータは含まれない

#### Scenario: スキーマエクスプローラーでの確認
- **WHEN** GraphQL Playgroundのスキーマエクスプローラーを開く
- **THEN** `Query.weatherAlerts(region: String, issuedAfter: String)` が表示される
- **AND** `WeatherAlert` タイプとそのフィールドが表示される

## MODIFIED Requirements

なし

## REMOVED Requirements

なし
# weather-alert-firestore Specification Delta

## ADDED Requirements

### Requirement: WeatherAlert Domain Model
システムは気象アラート本体データのドメインモデルを定義することをMUSTとする。

#### Scenario: WeatherAlert構造体の定義
- **WHEN** `internal/domain/weather_alert.go` を確認する
- **THEN** `WeatherAlert` 構造体が定義されている
- **AND** 以下のフィールドが含まれる:
  - `ID string` (Firestore: "id")
  - `Title string` (Firestore: "title")
  - `Description string` (Firestore: "description")
  - `RawData map[string]interface{}` (Firestore: "rawData")
  - `AffectedAreas []string` (Firestore: "affectedAreas")
  - `Recommendations []string` (Firestore: "recommendations")
- **AND** 各フィールドに適切な `firestore` タグが付与されている

#### Scenario: RawDataのJSON変換
- **WHEN** `WeatherAlert` をGraphQLレスポンスに変換する
- **THEN** `RawData map[string]interface{}` が JSON 文字列にシリアライズされる
- **AND** `json.Marshal()` を使用して変換される

### Requirement: WeatherAlert Repository Interface
システムは気象アラート本体データの取得を抽象化するリポジトリインターフェースを定義することをMUSTとする。

#### Scenario: WeatherAlertRepositoryインターフェース定義
- **WHEN** `internal/repository/weather_alert.go` を確認する
- **THEN** `WeatherAlertRepository` インターフェースが定義されている
- **AND** `GetByID(ctx context.Context, id string) (*WeatherAlert, error)` メソッドが含まれる
- **AND** `GetByIDs(ctx context.Context, ids []string) ([]*WeatherAlert, error)` メソッドが含まれる

#### Scenario: バッチ取得メソッドの定義
- **WHEN** `GetByIDs` メソッドのシグネチャを確認する
- **THEN** 複数のIDを受け取る `[]string` 引数が定義されている
- **AND** 複数の `WeatherAlert` を返す `[]*WeatherAlert` が返り値である

### Requirement: Firestore WeatherAlert Repository Implementation
システムはFirestoreを使用してWeatherAlertRepositoryを実装することをMUSTとする。

#### Scenario: FirestoreWeatherAlertRepository構造体の定義
- **WHEN** `internal/repository/firestore_weather_alert.go` を確認する
- **THEN** `FirestoreWeatherAlertRepository` 構造体が定義されている
- **AND** `client *firestore.Client` フィールドを持つ
- **AND** `WeatherAlertRepository` インターフェースを実装している

#### Scenario: GetByIDメソッドの実装
- **WHEN** `GetByID(ctx, "alert-001")` を呼び出す
- **THEN** `client.Collection("weatherAlerts").Doc("alert-001").Get(ctx)` が実行される
- **AND** Firestoreドキュメントが `WeatherAlert` 構造体にマッピングされる
- **AND** 成功時は `*WeatherAlert` が返される

#### Scenario: GetByIDメソッドのエラーハンドリング（ドキュメント未存在）
- **WHEN** 存在しないID `"nonexistent"` で `GetByID` を呼び出す
- **THEN** `nil` と `error` が返される
- **AND** エラーメッセージは "weather alert not found" を含む

#### Scenario: GetByIDsメソッドの実装（バッチ取得）
- **WHEN** `GetByIDs(ctx, []string{"alert-001", "alert-002", "alert-003"})` を呼び出す
- **THEN** Firestore の batch get 操作が実行される
- **AND** 各IDに対して `client.Collection("weatherAlerts").Doc(id).Get(ctx)` が実行される
- **AND** 取得成功したドキュメントのみが結果配列に含まれる

#### Scenario: GetByIDsメソッドの部分的な失敗ハンドリング
- **WHEN** 3つのIDのうち1つが存在しない場合
- **THEN** 存在する2つの `WeatherAlert` が返される
- **AND** 存在しないIDについてはログに警告が記録される
- **AND** エラーは返されない（部分的な成功を許容）

#### Scenario: GetByIDsメソッドの空配列ハンドリング
- **WHEN** `GetByIDs(ctx, []string{})` を呼び出す
- **THEN** 空の配列 `[]*WeatherAlert{}` が返される
- **AND** Firestoreクエリは実行されない

### Requirement: Firestore Collection Structure
システムは気象アラートデータを `weatherAlerts` コレクションに格納することをMUSTとする。

#### Scenario: コレクション名の定義
- **WHEN** Firestoreクライアントを使用してデータを取得する
- **THEN** コレクション名は `"weatherAlerts"` である
- **AND** ドキュメントIDは `weather_alert_metadata.id` と一致する

#### Scenario: ドキュメント構造の検証
- **WHEN** Firestoreドキュメントを確認する
- **THEN** 以下のフィールドが含まれる:
  - `id` (string)
  - `title` (string)
  - `description` (string)
  - `rawData` (map)
  - `affectedAreas` (array of strings)
  - `recommendations` (array of strings)

### Requirement: Repository Constructor and Dependency Injection
システムはWeatherAlertRepositoryの依存性注入を提供することをMUSTとする。

#### Scenario: NewFirestoreWeatherAlertRepositoryコンストラクタ
- **WHEN** `internal/repository/firestore_weather_alert.go` を確認する
- **THEN** `NewFirestoreWeatherAlertRepository(client *firestore.Client) *FirestoreWeatherAlertRepository` 関数が定義されている
- **AND** Firestoreクライアントを引数として受け取る
- **AND** 初期化された `FirestoreWeatherAlertRepository` を返す

#### Scenario: server.goでの依存性注入
- **WHEN** `server.go` を確認する
- **THEN** Firestoreクライアント初期化後に `NewFirestoreWeatherAlertRepository(firestoreClient)` が呼び出される
- **AND** リポジトリインスタンスがResolverに渡される

### Requirement: Error Handling and Logging
システムはFirestore操作のエラーを適切にハンドリングし、ログに記録することをMUSTとする。

#### Scenario: 接続エラーのハンドリング
- **WHEN** Firestore Emulatorが起動していない状態でクエリを実行する
- **THEN** Repositoryは接続エラーを返す
- **AND** エラーログが記録される

#### Scenario: 操作ログの記録
- **WHEN** `GetByID` または `GetByIDs` が実行される
- **THEN** 操作の開始と完了がログに記録される
- **AND** 取得したドキュメント数がログに記録される
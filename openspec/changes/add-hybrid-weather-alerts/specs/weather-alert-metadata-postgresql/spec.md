# weather-alert-metadata-postgresql Specification Delta

## ADDED Requirements

### Requirement: Weather Alert Metadata Schema
システムは気象アラートの検索用メタデータをPostgreSQLに格納するためのスキーマを定義することをMUSTとする。

#### Scenario: weather_alert_metadataテーブルの作成
- **WHEN** PostgreSQLマイグレーションスクリプトを実行する
- **THEN** `weather_alert_metadata` テーブルが作成される
- **AND** 以下のカラムが定義される:
  - `id` VARCHAR(255) PRIMARY KEY
  - `region` VARCHAR(100) NOT NULL
  - `severity` VARCHAR(50) NOT NULL
  - `issued_at` TIMESTAMP NOT NULL
  - `created_at` TIMESTAMP DEFAULT NOW()

#### Scenario: 検索用インデックスの作成
- **WHEN** PostgreSQLマイグレーションスクリプトを実行する
- **THEN** `region` カラムにインデックスが作成される
- **AND** `issued_at` カラムにインデックスが作成される
- **AND** `severity` カラムにインデックスが作成される

#### Scenario: マイグレーションファイルの配置
- **WHEN** プロジェクトディレクトリを確認する
- **THEN** `migrations/003_create_weather_alert_metadata.sql` ファイルが存在する
- **AND** ファイルには CREATE TABLE 文が含まれる
- **AND** ファイルには CREATE INDEX 文が含まれる

### Requirement: WeatherAlertMetadata Domain Model
システムは気象アラートメタデータのドメインモデルを定義することをMUSTとする。

#### Scenario: WeatherAlertMetadata構造体の定義
- **WHEN** `internal/domain/weather_alert_metadata.go` を確認する
- **THEN** `WeatherAlertMetadata` 構造体が定義されている
- **AND** 以下のフィールドが含まれる:
  - `ID string`
  - `Region string`
  - `Severity string`
  - `IssuedAt time.Time`
  - `CreatedAt time.Time`

#### Scenario: Severity定数の定義
- **WHEN** `internal/domain/weather_alert_metadata.go` を確認する
- **THEN** 重要度を表す定数が定義されている
- **AND** `SeverityInfo = "info"` が定義されている
- **AND** `SeverityWarning = "warning"` が定義されている
- **AND** `SeverityCritical = "critical"` が定義されている

### Requirement: WeatherAlertMetadata Repository Interface
システムは気象アラートメタデータの検索を抽象化するリポジトリインターフェースを定義することをMUSTとする。

#### Scenario: WeatherAlertMetadataRepositoryインターフェース定義
- **WHEN** `internal/repository/weather_alert_metadata.go` を確認する
- **THEN** `WeatherAlertMetadataRepository` インターフェースが定義されている
- **AND** `SearchIDs(ctx context.Context, filter MetadataFilter) ([]string, error)` メソッドが含まれる

#### Scenario: MetadataFilterの定義
- **WHEN** `internal/repository/weather_alert_metadata.go` を確認する
- **THEN** `MetadataFilter` 構造体が定義されている
- **AND** `Region *string` フィールドが含まれる
- **AND** `IssuedAfter *time.Time` フィールドが含まれる
- **AND** 各フィールドはポインタ型（オプショナル）である

### Requirement: PostgreSQL WeatherAlertMetadata Repository Implementation
システムはPostgreSQLを使用してWeatherAlertMetadataRepositoryを実装することをMUSTとする。

#### Scenario: PostgresWeatherAlertMetadataRepository構造体の定義
- **WHEN** `internal/repository/postgres_weather_alert_metadata.go` を確認する
- **THEN** `PostgresWeatherAlertMetadataRepository` 構造体が定義されている
- **AND** `db *sql.DB` フィールドを持つ
- **AND** `WeatherAlertMetadataRepository` インターフェースを実装している

#### Scenario: SearchIDsメソッドの実装（フィルタなし）
- **WHEN** `SearchIDs(ctx, MetadataFilter{})` を呼び出す
- **THEN** `SELECT id FROM weather_alert_metadata ORDER BY issued_at DESC` が実行される
- **AND** 全てのID配列が返される

#### Scenario: SearchIDsメソッドの実装（地域フィルタあり）
- **WHEN** `SearchIDs(ctx, MetadataFilter{Region: &"Tokyo"})` を呼び出す
- **THEN** `SELECT id FROM weather_alert_metadata WHERE region = $1 ORDER BY issued_at DESC` が実行される
- **AND** 地域が "Tokyo" のID配列が返される

#### Scenario: SearchIDsメソッドの実装（日時フィルタあり）
- **WHEN** `SearchIDs(ctx, MetadataFilter{IssuedAfter: &time})` を呼び出す
- **THEN** `SELECT id FROM weather_alert_metadata WHERE issued_at >= $1 ORDER BY issued_at DESC` が実行される
- **AND** 指定日時以降のID配列が返される

#### Scenario: SearchIDsメソッドの実装（複数フィルタあり）
- **WHEN** `SearchIDs(ctx, MetadataFilter{Region: &"Tokyo", IssuedAfter: &time})` を呼び出す
- **THEN** `SELECT id FROM weather_alert_metadata WHERE region = $1 AND issued_at >= $2 ORDER BY issued_at DESC` が実行される
- **AND** 両方の条件を満たすID配列が返される

#### Scenario: 空の結果のハンドリング
- **WHEN** 条件に合致するデータが存在しない
- **THEN** 空の配列 `[]string{}` が返される
- **AND** エラーは返されない

#### Scenario: データベースエラーのハンドリング
- **WHEN** PostgreSQL接続エラーが発生する
- **THEN** `nil` と `error` が返される
- **AND** エラーメッセージにコンテキストが含まれる

### Requirement: Repository Constructor and Dependency Injection
システムはWeatherAlertMetadataRepositoryの依存性注入を提供することをMUSTとする。

#### Scenario: NewPostgresWeatherAlertMetadataRepositoryコンストラクタ
- **WHEN** `internal/repository/postgres_weather_alert_metadata.go` を確認する
- **THEN** `NewPostgresWeatherAlertMetadataRepository(db *sql.DB) *PostgresWeatherAlertMetadataRepository` 関数が定義されている
- **AND** `sql.DB` を引数として受け取る
- **AND** 初期化された `PostgresWeatherAlertMetadataRepository` を返す

#### Scenario: main.goでの依存性注入
- **WHEN** `server.go` を確認する
- **THEN** PostgreSQL クライアント初期化後に `NewPostgresWeatherAlertMetadataRepository(db)` が呼び出される
- **AND** リポジトリインスタンスがResolverに渡される
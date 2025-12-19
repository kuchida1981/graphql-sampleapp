# Design: Hybrid Weather Alerts with PostgreSQL + Firestore

## Overview

本変更では、**PostgreSQLメタデータ検索 + Firestore本体データ取得**のハイブリッドパターンを実装します。気象アラートドメインを題材に、スキーマ変更に強い柔軟なアーキテクチャを示します。

## Background

別サービスで運用している課題:
- 気象データの元データスキーマが頻繁に変更される
- RDB スキーマに厳密に定義しているため、スキーマ変更の影響が大きい

解決策:
- **スキーマ変更に強いデータ**: Firestore (NoSQL) に格納
- **検索・集計・正規化が必要なデータ**: PostgreSQL に格納
- GraphQLクエリでは PostgreSQL で検索→ Firestore から詳細取得

## Domain Model: WeatherAlert

### PostgreSQL側のメタデータ (weather_alert_metadata)

検索・集計用のシンプルなメタデータのみ:

```sql
CREATE TABLE weather_alert_metadata (
    id VARCHAR(255) PRIMARY KEY,
    region VARCHAR(100) NOT NULL,         -- 地域 (例: "Tokyo", "Osaka")
    severity VARCHAR(50) NOT NULL,        -- 重要度 (例: "info", "warning", "critical")
    issued_at TIMESTAMP NOT NULL,         -- 発行日時
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_region ON weather_alert_metadata(region);
CREATE INDEX idx_issued_at ON weather_alert_metadata(issued_at);
CREATE INDEX idx_severity ON weather_alert_metadata(severity);
```

### Firestore側の本体データ (weatherAlerts collection)

詳細な観測データJSON (スキーマフリー):

```json
{
  "id": "alert-20251219-tokyo-001",
  "title": "Strong Wind Warning",
  "description": "Strong winds expected in Tokyo area",
  "rawData": {
    "temperature": { "value": 15.2, "unit": "celsius" },
    "windSpeed": { "value": 25.5, "unit": "m/s" },
    "precipitation": { "value": 0, "unit": "mm" },
    "pressure": { "value": 1013.2, "unit": "hPa" }
  },
  "affectedAreas": ["Chiyoda", "Minato", "Shibuya"],
  "recommendations": ["Stay indoors", "Secure loose objects"]
}
```

**ポイント**: `rawData` フィールドは気象庁のスキーマ変更に追従可能。PostgreSQLスキーマ変更は不要。

## Architecture

### Data Flow

```
GraphQL Query (region: "Tokyo", issuedAfter: "2025-12-19T00:00:00Z")
    ↓
PostgreSQL Query
    SELECT id FROM weather_alert_metadata
    WHERE region = 'Tokyo' AND issued_at >= '2025-12-19T00:00:00Z'
    → Result: ["alert-001", "alert-002", "alert-003"]
    ↓
Firestore Batch Get
    firestore.Collection("weatherAlerts").Doc("alert-001").Get()
    firestore.Collection("weatherAlerts").Doc("alert-002").Get()
    firestore.Collection("weatherAlerts").Doc("alert-003").Get()
    ↓
GraphQL Response (merged data)
    [
      { id: "alert-001", region: "Tokyo", severity: "warning", rawData: {...} },
      ...
    ]
```

### Component Structure

```
graph/schema.graphqls
  ├─ WeatherAlert type
  └─ Query.weatherAlerts(region, issuedAfter)

graph/schema.resolvers.go
  └─ weatherAlerts resolver
      ├─ 1. PostgreSQL で ID リストを取得
      └─ 2. Firestore で詳細データを取得

internal/repository/
  ├─ weather_alert_metadata.go (PostgreSQL Repository Interface)
  ├─ postgres_weather_alert_metadata.go (PostgreSQL Implementation)
  ├─ weather_alert.go (Firestore Repository Interface)
  └─ firestore_weather_alert.go (Firestore Implementation)

internal/domain/
  ├─ weather_alert_metadata.go (PostgreSQL domain model)
  └─ weather_alert.go (Firestore domain model)
```

## Repository Pattern

### PostgreSQL Repository

```go
type WeatherAlertMetadataRepository interface {
    SearchIDs(ctx context.Context, filter MetadataFilter) ([]string, error)
}

type MetadataFilter struct {
    Region      *string
    IssuedAfter *time.Time
}

type PostgresWeatherAlertMetadataRepository struct {
    db *sql.DB
}
```

### Firestore Repository

```go
type WeatherAlertRepository interface {
    GetByID(ctx context.Context, id string) (*WeatherAlert, error)
    GetByIDs(ctx context.Context, ids []string) ([]*WeatherAlert, error)
}

type FirestoreWeatherAlertRepository struct {
    client *firestore.Client
}
```

## GraphQL Schema Design

```graphql
type WeatherAlert {
  id: ID!
  region: String!           # from PostgreSQL metadata
  severity: String!         # from PostgreSQL metadata
  issuedAt: String!         # from PostgreSQL metadata
  title: String!            # from Firestore
  description: String!      # from Firestore
  rawData: String!          # from Firestore (JSON string)
}

type Query {
  weatherAlerts(
    region: String
    issuedAfter: String
  ): [WeatherAlert!]!
}
```

## Trade-offs and Considerations

### Pros
- ✅ Firestore スキーマ変更に強い（rawData フィールドはJSON文字列として柔軟に格納）
- ✅ PostgreSQL で効率的な検索・集計が可能
- ✅ 既存の User (PostgreSQL) と Message (Firestore) パターンを踏襲

### Cons
- ⚠️ 2つのデータストアを同期する必要がある（トランザクション保証なし）
- ⚠️ N+1問題のリスク（複数IDをFirestoreから取得）
- ⚠️ データ整合性の課題（PostgreSQLにIDがあるがFirestoreにデータがない、など）

### Mitigation
- **同期**: データ投入時に PostgreSQL → Firestore を順次実行（エラー時はロールバック）
- **N+1**: Firestore の batch get を使用（最大500件まで）
- **整合性**: Firestore取得失敗時は nil を返し、GraphQL レスポンスから除外

## Migration and Seeding

### PostgreSQL Migration

```sql
-- migrations/003_create_weather_alert_metadata.sql
CREATE TABLE weather_alert_metadata (
    id VARCHAR(255) PRIMARY KEY,
    region VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_region ON weather_alert_metadata(region);
CREATE INDEX idx_issued_at ON weather_alert_metadata(issued_at);
CREATE INDEX idx_severity ON weather_alert_metadata(severity);
```

### Seed Script

`scripts/seed-weather-alerts.go`:
1. PostgreSQL に 10件のメタデータを投入
2. Firestore に対応する 10件の詳細データを投入

## Testing Strategy

### Unit Tests
- Repository層のモック化
- PostgreSQL repository: filter条件のテスト
- Firestore repository: batch get のテスト

### Integration Tests
- PostgreSQL + Firestore の結合テスト
- GraphQL Playground でクエリ実行テスト

## Open Questions

1. ~~どのフィルタ条件を実装するか?~~ → 地域フィルタ、日時範囲検索（ユーザー確認済み）
2. Firestore データが見つからない場合の挙動 → nullを返すか、エラーとするか → **nullを返し、GraphQLレスポンスから除外**
3. ページネーションは必要か? → 初期実装では不要（シンプルさ優先）

## References

- 既存実装: `add-postgresql-user-integration`, `add-firestore-docker-integration`
- GraphQL スキーマ: `graph/schema.graphqls`
- Repository パターン: `internal/repository/*.go`
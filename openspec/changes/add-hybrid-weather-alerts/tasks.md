# Implementation Tasks: add-hybrid-weather-alerts

## Phase 1: PostgreSQL Metadata Layer

### Task 1.1: Create weather_alert_metadata schema migration
- Create `migrations/003_create_weather_alert_metadata.sql`
- Define `weather_alert_metadata` table with columns: id, region, severity, issued_at, created_at
- Add indexes on region, issued_at, severity
- **Validation**: Run migration and verify table exists with `\d weather_alert_metadata`

### Task 1.2: Define WeatherAlertMetadata domain model
- Create `internal/domain/weather_alert_metadata.go`
- Define `WeatherAlertMetadata` struct with fields: ID, Region, Severity, IssuedAt, CreatedAt
- Define severity constants: SeverityInfo, SeverityWarning, SeverityCritical
- **Validation**: Build succeeds with `go build`

### Task 1.3: Define WeatherAlertMetadataRepository interface
- Create `internal/repository/weather_alert_metadata.go`
- Define `WeatherAlertMetadataRepository` interface with `SearchIDs` method
- Define `MetadataFilter` struct with optional Region and IssuedAfter fields
- **Validation**: Build succeeds with `go build`

### Task 1.4: Implement PostgresWeatherAlertMetadataRepository
- Create `internal/repository/postgres_weather_alert_metadata.go`
- Implement `PostgresWeatherAlertMetadataRepository` struct with `db *sql.DB`
- Implement `SearchIDs` method with dynamic WHERE clause construction
- Add constructor `NewPostgresWeatherAlertMetadataRepository(db *sql.DB)`
- **Validation**: Build succeeds, run unit test for filter logic
- **Dependencies**: Task 1.3 must be completed

## Phase 2: Firestore Data Layer

### Task 2.1: Define WeatherAlert domain model
- Create `internal/domain/weather_alert.go`
- Define `WeatherAlert` struct with fields: ID, Title, Description, RawData, AffectedAreas, Recommendations
- Add `firestore` tags to all fields
- **Validation**: Build succeeds with `go build`

### Task 2.2: Define WeatherAlertRepository interface
- Create `internal/repository/weather_alert.go`
- Define `WeatherAlertRepository` interface with `GetByID` and `GetByIDs` methods
- **Validation**: Build succeeds with `go build`

### Task 2.3: Implement FirestoreWeatherAlertRepository
- Create `internal/repository/firestore_weather_alert.go`
- Implement `FirestoreWeatherAlertRepository` struct with `client *firestore.Client`
- Implement `GetByID` method to fetch single document from `weatherAlerts` collection
- Implement `GetByIDs` method with batch get logic (iterate and collect results)
- Add constructor `NewFirestoreWeatherAlertRepository(client *firestore.Client)`
- Add error handling for missing documents (log warning, continue)
- **Validation**: Build succeeds, run unit test with Firestore emulator
- **Dependencies**: Task 2.2 must be completed

## Phase 3: GraphQL Integration

### Task 3.1: Update GraphQL schema
- Edit `graph/schema.graphqls`
- Add `WeatherAlert` type with fields: id, region, severity, issuedAt, title, description, rawData, affectedAreas, recommendations
- Add `weatherAlerts(region: String, issuedAfter: String): [WeatherAlert!]!` query to Query type
- **Validation**: Run `go run github.com/99designs/gqlgen generate` and check for errors

### Task 3.2: Update Resolver struct and constructor
- Edit `graph/resolver.go`
- Add `weatherAlertMetadataRepo` and `weatherAlertRepo` fields to `Resolver` struct
- Update `NewResolver` constructor to accept both repositories as arguments
- **Validation**: Build succeeds with `go build`
- **Dependencies**: Task 1.3, Task 2.2 must be completed

### Task 3.3: Implement weatherAlerts resolver
- Edit `graph/schema.resolvers.go`
- Implement `WeatherAlerts` resolver method
- Step 1: Build `MetadataFilter` from query arguments (region, issuedAfter)
- Step 2: Call `weatherAlertMetadataRepo.SearchIDs(ctx, filter)` to get ID list
- Step 3: Call `weatherAlertRepo.GetByIDs(ctx, ids)` to get Firestore data
- Step 4: Merge PostgreSQL metadata with Firestore data into `model.WeatherAlert`
- Step 5: Marshal `rawData` map to JSON string
- Add error handling for PostgreSQL errors, Firestore errors, and date parsing errors
- **Validation**: Build succeeds, test with GraphQL Playground
- **Dependencies**: Task 1.4, Task 2.3, Task 3.1, Task 3.2 must be completed

### Task 3.4: Wire repositories in server.go
- Edit `server.go`
- After PostgreSQL client initialization, create `NewPostgresWeatherAlertMetadataRepository(db)`
- After Firestore client initialization, create `NewFirestoreWeatherAlertRepository(firestoreClient)`
- Pass both repositories to `graph.NewResolver(...)`
- **Validation**: Run `go run server.go` and verify no initialization errors
- **Dependencies**: Task 3.2 must be completed

## Phase 4: Data Seeding and Testing

### Task 4.1: Run PostgreSQL migration
- Run migration script: `docker compose exec app psql $DATABASE_URL -f migrations/003_create_weather_alert_metadata.sql` (or appropriate method)
- Verify table creation with `\d weather_alert_metadata`
- **Validation**: Table exists with correct schema
- **Dependencies**: Task 1.1 must be completed

### Task 4.2: Create seed script for weather alerts
- Create `scripts/seed-weather-alerts.go`
- Seed 10 weather alert metadata records to PostgreSQL with varied regions (Tokyo, Osaka, Kyoto) and dates
- Seed corresponding 10 weather alert documents to Firestore `weatherAlerts` collection
- Include sample `rawData` with temperature, windSpeed, precipitation, pressure
- **Validation**: Run script and verify data in both databases
- **Dependencies**: Task 1.1, Task 2.1 must be completed

### Task 4.3: Test weatherAlerts query in GraphQL Playground
- Start server with `go run server.go`
- Open GraphQL Playground at `http://localhost:8080/`
- Test query: `{ weatherAlerts { id region severity title } }`
- Test query with region filter: `{ weatherAlerts(region: "Tokyo") { id region } }`
- Test query with date filter: `{ weatherAlerts(issuedAfter: "2025-12-19T00:00:00Z") { id issuedAt } }`
- Verify correct data is returned
- **Validation**: All queries return expected results
- **Dependencies**: Task 3.3, Task 4.2 must be completed

## Phase 5: Code Quality and Documentation

### Task 5.1: Run gofmt and goimports
- Run `gofmt -w .` to format all Go code
- Run `goimports -w .` to organize imports
- **Validation**: No changes after re-running commands

### Task 5.2: Add README documentation
- Update `README.md` with new WeatherAlert feature description
- Document GraphQL query examples
- Document data flow (PostgreSQL → Firestore)
- **Validation**: Documentation is clear and accurate

### Task 5.3: Verify all tests pass
- Run `go test ./...` to execute all tests
- Fix any failing tests
- **Validation**: All tests pass
- **Dependencies**: All implementation tasks must be completed

## Task Dependencies Summary

```
Phase 1 (PostgreSQL)          Phase 2 (Firestore)
Task 1.1 → Task 1.2           Task 2.1
       ↓                             ↓
Task 1.3 → Task 1.4           Task 2.2 → Task 2.3
       ↓                             ↓
       └─────────┬───────────────────┘
                 ↓
         Phase 3 (GraphQL)
         Task 3.1 → Task 3.2 → Task 3.3 → Task 3.4
                                 ↓
                         Phase 4 (Seeding)
                         Task 4.1 → Task 4.2 → Task 4.3
                                 ↓
                         Phase 5 (Quality)
                         Task 5.1 → Task 5.2 → Task 5.3
```

## Notes
- Tasks within the same phase can be parallelized if no explicit dependency is listed
- Each task should be completed and validated before moving to dependent tasks
- All code should follow Go conventions and project code style guidelines
- Use existing User (PostgreSQL) and Message (Firestore) implementations as reference
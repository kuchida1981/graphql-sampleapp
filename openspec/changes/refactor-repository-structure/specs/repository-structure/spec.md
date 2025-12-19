# Spec Deltas for Repository Structure

## MODIFIED Requirements

### capability: `repository-structure`

### Requirement: The repository directory structure must be organized by data store.
The repository directory structure MUST be organized into subdirectories based on the data store implementation.

#### Scenario: Directory Layout
- **Given** the project's directory structure.
- **When** a developer inspects the `internal/repository` directory.
- **Then** it MUST contain the generic repository interface files (`user.go`, `message.go`, etc.).
- **And** it MUST contain a `postgres` subdirectory for PostgreSQL-specific implementations.
- **And** it MUST contain a `firestore` subdirectory for Firestore-specific implementations.

#### Scenario: PostgreSQL Implementation Location
- **Given** the `internal/repository/postgres` directory.
- **When** a developer inspects its contents.
- **Then** it MUST contain the implementations for the `UserRepository` and `WeatherAlertMetadataRepository`.

#### Scenario: Firestore Implementation Location
- **Given** the `internal/repository/firestore` directory.
- **When** a developer inspects its contents.
- **Then** it MUST contain the implementations for the `MessageRepository` and `WeatherAlertRepository`.

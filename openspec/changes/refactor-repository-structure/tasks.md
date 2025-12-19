# Tasks: Refactor Repository Structure

1.  [x] Create new directories: `internal/repository/postgres` and `internal/repository/firestore`.
2.  [x] Move PostgreSQL-specific files (`postgres_*.go`) into the `internal/repository/postgres` directory and rename them (e.g., `postgres_user.go` -> `user.go`).
3.  [x] Move Firestore-specific files (`firestore_*.go`) into the `internal/repository/firestore` directory and rename them (e.g., `firestore_message.go` -> `message.go`).
4.  [x] Update all import paths in the codebase that reference the moved repository implementations. This will primarily affect `server.go` where repositories are instantiated.
5.  [x] Run `go mod tidy` to clean up dependencies.
6.  [x] Run all tests with `go test ./...` to ensure the application still functions correctly after the refactoring.
7.  [x] Remove the old `firestore_*.go` and `postgres_*.go` files from `internal/repository`.

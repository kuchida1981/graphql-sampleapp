//go:build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jxpress/graphql-sampleapp/internal/postgres"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://graphql_user:graphql_pass@localhost:5432/graphql_db?sslmode=disable"
		log.Printf("DATABASE_URL not set, using default: %s", databaseURL)
	}

	db, err := postgres.NewClient(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	users := []struct {
		id        string
		name      string
		email     string
		createdAt time.Time
	}{
		{"user1", "Alice Smith", "alice@example.com", time.Now().Add(-48 * time.Hour)},
		{"user2", "Bob Johnson", "bob@example.com", time.Now().Add(-24 * time.Hour)},
		{"user3", "Charlie Brown", "charlie@example.com", time.Now().Add(-12 * time.Hour)},
		{"user4", "Diana Prince", "diana@example.com", time.Now().Add(-6 * time.Hour)},
		{"user5", "Eve Adams", "eve@example.com", time.Now()},
	}

	query := `
		INSERT INTO users (id, name, email, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    email = EXCLUDED.email,
		    created_at = EXCLUDED.created_at
	`

	for _, user := range users {
		_, err := db.ExecContext(ctx, query, user.id, user.name, user.email, user.createdAt)
		if err != nil {
			log.Fatalf("Failed to insert user %s: %v", user.id, err)
		}
		log.Printf("Seeded user: %s (%s)", user.name, user.email)
	}

	log.Println("Successfully seeded PostgreSQL database with sample users")
}

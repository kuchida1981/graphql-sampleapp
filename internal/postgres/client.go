package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewClient(ctx context.Context, connStr string) (*sql.DB, error) {
	if connStr == "" {
		connStr = os.Getenv("DATABASE_URL")
		if connStr == "" {
			return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
		}
	}

	log.Printf("Connecting to PostgreSQL: %s", maskPassword(connStr))

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}

func maskPassword(connStr string) string {
	return "postgres://***:***@<host>/<db>"
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/kuchida1981/graphql-sampleapp/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	log.Println("PostgresUserRepository: Listing all users")

	query := "SELECT id, name, email, created_at FROM users ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("PostgresUserRepository: Failed to query users: %v", err)
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			log.Printf("PostgresUserRepository: Failed to scan user: %v", err)
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("PostgresUserRepository: Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("PostgresUserRepository: Found %d users", len(users))
	return users, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	log.Printf("PostgresUserRepository: Getting user by ID: %s", id)

	query := "SELECT id, name, email, created_at FROM users WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("PostgresUserRepository: User not found: %s", id)
			return nil, fmt.Errorf("user not found: %s", id)
		}
		log.Printf("PostgresUserRepository: Failed to scan user: %v", err)
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	log.Printf("PostgresUserRepository: Found user: %s", user.ID)
	return &user, nil
}

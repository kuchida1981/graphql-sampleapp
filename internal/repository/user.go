package repository

import (
	"context"

	"github.com/kuchida1981/graphql-sampleapp/internal/domain"
)

type UserRepository interface {
	List(ctx context.Context) ([]*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

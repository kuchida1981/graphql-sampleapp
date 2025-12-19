package repository

import (
	"context"

	"github.com/kuchida1981/graphql-sampleapp/internal/domain"
)

type MessageRepository interface {
	List(ctx context.Context) ([]*domain.Message, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
}

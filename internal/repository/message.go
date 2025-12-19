package repository

import (
	"context"

	"github.com/jxpress/graphql-sampleapp/internal/domain"
)

type MessageRepository interface {
	List(ctx context.Context) ([]*domain.Message, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
}

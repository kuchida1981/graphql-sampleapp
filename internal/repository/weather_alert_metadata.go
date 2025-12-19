package repository

import (
	"context"
	"time"

	"github.com/kuchida1981/graphql-sampleapp/internal/domain"
)

type MetadataFilter struct {
	Region      *string
	IssuedAfter *time.Time
}

type WeatherAlertMetadataRepository interface {
	SearchIDs(ctx context.Context, filter MetadataFilter) ([]string, error)
	Search(ctx context.Context, filter MetadataFilter) ([]*domain.WeatherAlertMetadata, error)
}

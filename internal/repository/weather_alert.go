package repository

import (
	"context"

	"github.com/jxpress/graphql-sampleapp/internal/domain"
)

type WeatherAlertRepository interface {
	GetByID(ctx context.Context, id string) (*domain.WeatherAlert, error)
	GetByIDs(ctx context.Context, ids []string) ([]*domain.WeatherAlert, error)
}

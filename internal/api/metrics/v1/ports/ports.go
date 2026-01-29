package ports

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/entities"
)

type MetricsRepository interface {
	Upload(ctx context.Context, metrics entities.Metrics) error
	Download(ctx context.Context) ([]entities.Metrics, error)
}

type Store struct {
	Metrics MetricsRepository
}

type StoreProvider interface {
	Provide(ctx context.Context) Store
}

package ram

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/ports"
)

type StoreProvider struct {
	repo ports.MetricsRepository
}

func NewStoreProvider() *StoreProvider {
	return &StoreProvider{repo: NewMetricsRepository()}
}

func (p *StoreProvider) Provide(_ context.Context) ports.Store {
	return ports.Store{Metrics: p.repo}
}

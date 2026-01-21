package ram

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/entities"
)

type MetricsRepository struct {
	data []entities.Metrics
}

func NewMetricsRepository() *MetricsRepository {
	return &MetricsRepository{data: make([]entities.Metrics, 0)}
}

func (m *MetricsRepository) Upload(_ context.Context, metrics entities.Metrics) error {
	m.data = append(m.data, metrics)

	return nil
}

func (m *MetricsRepository) Download(_ context.Context) ([]entities.Metrics, error) {
	return m.data, nil
}

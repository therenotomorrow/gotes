package ram

import (
	"context"
	"sync"

	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/entities"
)

type MetricsRepository struct {
	data  []entities.Metrics
	mutex sync.RWMutex
}

func NewMetricsRepository() *MetricsRepository {
	return &MetricsRepository{
		mutex: sync.RWMutex{},
		data:  make([]entities.Metrics, 0),
	}
}

func (m *MetricsRepository) Upload(_ context.Context, metrics entities.Metrics) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = append(m.data, metrics)

	return nil
}

func (m *MetricsRepository) Download(_ context.Context) ([]entities.Metrics, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.data, nil
}

package usecases

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/entities"
	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/ports"
)

const (
	percents = 100
)

type UseCases struct {
	store ports.Store
}

func NewCases(store ports.Store) *UseCases {
	return &UseCases{store: store}
}

func (use *UseCases) UploadMetrics(ctx context.Context, metrics entities.Metrics) error {
	return use.store.Metrics.Upload(ctx, metrics)
}

func (use *UseCases) CollectMetrics(ctx context.Context) (entities.Statistics, error) {
	var stat entities.Statistics

	metrics, err := use.store.Metrics.Download(ctx)
	if err != nil {
		return stat, err
	}

	errs := int64(0)

	for _, metric := range metrics {
		stat.Total += metric.Requests
		errs += metric.Errors
	}

	stat.ErrorRate = (float64(errs) / float64(stat.Total)) * percents

	return stat, nil
}

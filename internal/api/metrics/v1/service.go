package v1

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api"
	adapters "github.com/therenotomorrow/gotes/internal/api/metrics/v1/adapters/ram"
	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/entities"
	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/ports"
	"github.com/therenotomorrow/gotes/internal/api/metrics/v1/usecases"
	pb "github.com/therenotomorrow/gotes/pkg/api/metrics/v1"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"google.golang.org/grpc"
)

const (
	ErrStream ex.Error = "stream error"
)

type MetricsService struct {
	pb.UnimplementedMetricsServiceServer

	handle api.ErrorHandlerFunc
	tracer *trace.Tracer
	cases  *usecases.UseCases
}

func NewService(logger *slog.Logger) *MetricsService {
	provider := adapters.NewStoreProvider()

	return NewServiceWithProvider(provider, logger)
}

func NewServiceWithProvider(provider ports.StoreProvider, logger *slog.Logger) *MetricsService {
	store := provider.Provide(context.Background())

	return &MetricsService{
		UnimplementedMetricsServiceServer: pb.UnimplementedMetricsServiceServer{},
		handle:                            api.ErrorHandler(NewErrorMarshaler()),
		tracer:                            trace.Service("metrics.v1", logger),
		cases:                             usecases.NewCases(store),
	}
}

func (svc *MetricsService) UploadMetrics(
	stream grpc.ClientStreamingServer[pb.UploadMetricsRequest, pb.UploadMetricsResponse],
) error {
	ctx := stream.Context()

	for {
		req, err := stream.Recv()

		switch {
		case errors.Is(err, io.EOF):
			return svc.uploadMetricsSendAndClose(stream)
		case err != nil:
			return svc.handle(ErrStream.Because(err))
		}

		err = svc.cases.UploadMetrics(ctx, entities.Metrics{
			Requests: req.GetRequests(),
			Errors:   req.GetErrors(),
		})
		if err != nil {
			return ErrStream.Because(err)
		}
	}
}

func (svc *MetricsService) uploadMetricsSendAndClose(
	stream grpc.ClientStreamingServer[pb.UploadMetricsRequest, pb.UploadMetricsResponse],
) error {
	ctx := stream.Context()

	stat, err := svc.cases.CollectMetrics(ctx)
	if err != nil {
		return svc.handle(ErrStream.Because(err))
	}

	err = stream.SendAndClose(&pb.UploadMetricsResponse{
		Total:     stat.Total,
		ErrorRate: stat.ErrorRate,
	})
	if err != nil {
		return svc.handle(ErrStream.Because(err))
	}

	return nil
}

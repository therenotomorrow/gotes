package main

import (
	"context"
	"math/rand/v2"
	"sync"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/config"
	pb "github.com/therenotomorrow/gotes/pkg/api/metrics/v1"
	"github.com/therenotomorrow/gotes/pkg/client"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
)

func main() {
	var (
		ctx = context.Background()
		cfg = config.MustNew()
		log = trace.Logger(trace.TEXT, cfg.Debug)
		cli = client.MustNew(client.Config{
			Address: cfg.Server.Address,
			Secure:  false,
		})
		wait = sync.WaitGroup{}
	)

	defer cli.Close()

	var requests int64 = 100

	for range 5 {
		wait.Go(func() {
			stream, err := cli.UploadMetrics(ctx)
			ex.Panic(err)

			for range 10 {
				err = stream.Send(&pb.UploadMetricsRequest{
					Requests: requests,
					Errors:   rand.Int64N(requests), //nolint:gosec // allowed for simplicity
				})
				ex.Panic(err)
			}

			resp, err := stream.CloseAndRecv()
			ex.Panic(err)

			log.Info("UploadMetrics response", "resp", resp)
		})
	}

	wait.Wait()

	stream, err := cli.UploadMetrics(ctx)
	ex.Panic(err)

	total, err := stream.CloseAndRecv()
	ex.Panic(err)

	log.Info("UploadMetrics total", "total", total)
}

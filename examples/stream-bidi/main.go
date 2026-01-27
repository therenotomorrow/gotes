package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/config"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
	pb "github.com/therenotomorrow/gotes/pkg/api/chat/v1"
	"github.com/therenotomorrow/gotes/pkg/client"
	"github.com/therenotomorrow/gotes/pkg/services/generate"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"google.golang.org/grpc"
)

const (
	halfSecond = time.Second / 2
)

func recv(stream grpc.BidiStreamingClient[pb.DispatchRequest, pb.DispatchResponse], log *slog.Logger) {
	for {
		resp, err := stream.Recv()

		switch {
		case errors.Is(err, io.EOF):
			return
		case err != nil:
			log.Error("recv error", "error", err)
			ex.Panic(err)
		}

		switch payload := resp.GetPayload().(type) {
		case *pb.DispatchResponse_Status:
			log.Info("recv status", "status", payload.Status)
		case *pb.DispatchResponse_Message:
			log.Info("recv message", "message", payload.Message)
		}
	}
}

func send(stream grpc.BidiStreamingClient[pb.DispatchRequest, pb.DispatchResponse]) {
	messages := []string{
		"hello",
		"error", // simulate business error
		"how are you?",
		"oh", // simulate validation error
		"ping",
		"error", // simulate business error
		"good day",
	}

	for i, text := range messages {
		corrID := uuid.New()

		if i%4 == 0 {
			// simulate generation of correlation ID on the server side
			corrID = uuid.UUID{}
		}

		req := &pb.DispatchRequest{Message: &pb.Message{
			Header: &pb.Header{CorrelationId: corrID.Value()},
			Text:   text,
		}}

		err := stream.Send(req)
		ex.Panic(err)

		time.Sleep(halfSecond)
	}

	err := stream.CloseSend()
	ex.Panic(err)
}

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

	uuid.SetGenerator(generate.NewUUID())

	defer cli.Close()

	stream, err := cli.Dispatch(ctx)
	ex.Panic(err)

	wait.Go(func() {
		send(stream)
		log.Info("send finished")
	})

	wait.Go(func() {
		recv(stream, log)
		log.Info("recv finished")
	})

	wait.Wait()
}

package v1

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"time"

	"buf.build/go/protovalidate"
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api"
	adapters "github.com/therenotomorrow/gotes/internal/api/chat/v1/adapters/ram"
	"github.com/therenotomorrow/gotes/internal/api/chat/v1/entities"
	"github.com/therenotomorrow/gotes/internal/api/chat/v1/ports"
	pb "github.com/therenotomorrow/gotes/pkg/api/chat/v1"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

const (
	numParts = 2
	deadline = 4 * time.Second

	ErrChat   ex.Error = "chat error"
	ErrFinish ex.Error = "finish error"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer

	handle    api.ErrorHandlerFunc
	store     ports.Store
	tracer    *trace.Tracer
	validator protovalidate.Validator
}

func NewService(validator protovalidate.Validator, logger *slog.Logger) *ChatService {
	provider := adapters.NewStoreProvider()

	return NewServiceWithProvider(validator, provider, logger)
}

func NewServiceWithProvider(
	validator protovalidate.Validator,
	provider ports.StoreProvider,
	logger *slog.Logger,
) *ChatService {
	store := provider.Provide(context.Background())

	return &ChatService{
		UnimplementedChatServiceServer: pb.UnimplementedChatServiceServer{},
		handle:                         api.ErrorHandler(NewErrorMarshaler()),
		tracer:                         trace.Service("chat.v1", logger),
		validator:                      validator,
		store:                          store,
	}
}

func (svc *ChatService) Dispatch(stream grpc.BidiStreamingServer[pb.DispatchRequest, pb.DispatchResponse]) error {
	group, ctx := errgroup.WithContext(stream.Context())

	group.Go(func() error {
		err := svc.dispatchRecv(ctx, stream)
		if errors.Is(err, ErrFinish) {
			return nil
		}

		return err
	})

	group.Go(func() error {
		err := svc.dispatchSend(ctx, stream)
		if errors.Is(err, ErrFinish) {
			return nil
		}

		return err
	})

	err := group.Wait()
	if err != nil {
		return svc.handle(ErrChat.Because(err))
	}

	return nil
}

func (svc *ChatService) dispatchRecv(
	ctx context.Context,
	stream grpc.BidiStreamingServer[pb.DispatchRequest, pb.DispatchResponse],
) error {
	for {
		req, err := stream.Recv()

		switch {
		case errors.Is(err, io.EOF):
			return ErrFinish
		case err != nil:
			return err
		}

		err = svc.validator.Validate(req)
		if err != nil {
			err = svc.dispatchRecvStatus(stream, err)
			if err != nil {
				return err
			}

			continue
		}

		message, err := entities.NewMessage(
			req.GetMessage().GetText(),
			req.GetMessage().GetHeader().GetCorrelationId(),
		)
		if err != nil {
			err = svc.dispatchRecvStatus(stream, err)
			if err != nil {
				return err
			}

			continue
		}

		err = svc.store.Messages.SaveMessage(ctx, message)
		if err != nil {
			return err
		}
	}
}

func (svc *ChatService) dispatchSend(
	ctx context.Context,
	stream grpc.BidiStreamingServer[pb.DispatchRequest, pb.DispatchResponse],
) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(deadline)
	defer timer.Stop()

	for {
		messages, err := svc.store.Messages.Outbox(ctx)
		if err != nil {
			return svc.handle(ErrChat.Because(err))
		}

		select {
		case <-timer.C:
			return ErrFinish
		default:
			time.Sleep(time.Second)
		}

		for _, message := range messages {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-ticker.C:
				if message.Text == "error" {
					err = svc.dispatchSendStatus(stream, message)
				} else {
					err = svc.dispatchSendMessage(stream, message)
				}

				if err != nil {
					return err
				}

				err = svc.store.Messages.DeleteMessage(ctx, message)
				ex.Skip(err)

				timer.Reset(deadline)
			}
		}
	}
}

func (svc *ChatService) dispatchRecvStatus(
	stream grpc.BidiStreamingServer[pb.DispatchRequest, pb.DispatchResponse],
	err error,
) error {
	if err == nil {
		return nil
	}

	st := status.New(codes.Unknown, err.Error())
	details := []protoadapt.MessageV1{&typespb.Error{
		Code:   typespb.ErrorCode_ERROR_CODE_UNKNOWN,
		Reason: err.Error(),
	}}

	verr := new(protovalidate.ValidationError)
	if errors.As(err, &verr) {
		st = status.New(codes.InvalidArgument, err.Error())
		details = make([]protoadapt.MessageV1, 0)

		for _, violation := range verr.Violations {
			parts := strings.SplitN(violation.String(), ": ", numParts)

			name := parts[0]
			code := typespb.ErrorCode_ERROR_CODE_UNKNOWN

			if name == "text" {
				code = typespb.ErrorCode_ERROR_CODE_INVALID_TEXT
			}

			details = append(details, &typespb.Error{Code: code, Reason: parts[1]})
		}
	}

	st, err = st.WithDetails(details...)

	ex.Skip(err)

	return stream.Send(&pb.DispatchResponse{Payload: MarshalStatus(st)})
}

func (svc *ChatService) dispatchSendStatus(
	stream grpc.BidiStreamingServer[pb.DispatchRequest, pb.DispatchResponse],
	message *entities.Message,
) error {
	st, err := status.New(codes.Unknown, message.Text).WithDetails(&typespb.Error{
		Code:   typespb.ErrorCode_ERROR_CODE_BUSINESS,
		Reason: "some business issues",
	})

	ex.Skip(err)

	return stream.Send(&pb.DispatchResponse{Payload: MarshalStatus(st)})
}

func (svc *ChatService) dispatchSendMessage(
	stream grpc.BidiStreamingServer[pb.DispatchRequest, pb.DispatchResponse],
	message *entities.Message,
) error {
	return stream.Send(&pb.DispatchResponse{
		Payload: MarshalMessage(&entities.Message{
			Header: message.Header,
			Text:   "processed: " + message.Text,
		}),
	})
}

package v1

import (
	"context"
	"log/slog"

	"github.com/therenotomorrow/gotes/internal/api"
	adapters "github.com/therenotomorrow/gotes/internal/api/users/v1/adapters/postgres"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/ports"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	pb "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
)

type UsersService struct {
	pb.UnimplementedUsersServiceServer

	handle api.ErrorHandlerFunc
	tracer *trace.Tracer
	cases  *usecases.UseCases
}

func New(db postgres.Database, logger *slog.Logger) *UsersService {
	provider := adapters.NewStoreProvider(db)

	return NewService(adapters.NewUnitOfWork(db, provider), logger)
}

func NewService(uow ports.UnitOfWork, logger *slog.Logger) *UsersService {
	return &UsersService{
		UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
		handle:                          api.ErrorHandler(NewErrorMarshaler()),
		tracer:                          trace.Service("users.v1", logger),
		cases:                           usecases.NewCases(uow),
	}
}

func (svc *UsersService) RegisterUser(
	ctx context.Context,
	request *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {
	user, err := svc.cases.RegisterUser(ctx, &usecases.RegisterUserInput{
		Name:     request.GetName(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, svc.handle(err)
	}

	return &pb.RegisterUserResponse{User: MarshalUser(user)}, nil
}

func (svc *UsersService) RefreshToken(
	ctx context.Context,
	request *pb.RefreshTokenRequest,
) (*pb.RefreshTokenResponse, error) {
	user, err := svc.cases.RefreshToken(ctx, &usecases.RefreshTokenInput{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, svc.handle(err)
	}

	return &pb.RefreshTokenResponse{User: MarshalUser(user)}, nil
}

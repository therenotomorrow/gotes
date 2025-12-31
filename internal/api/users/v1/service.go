package v1

import (
	"context"
	"log/slog"

	"github.com/therenotomorrow/gotes/internal/api"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/adapters"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	pb "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
)

type UsersService struct {
	pb.UnimplementedUsersServiceServer

	tracer *trace.Tracer
	cases  *UseCases
}

func Service(database postgres.Database, logger *slog.Logger) *UsersService {
	var (
		uow    = adapters.NewUnitOfWork(database)
		tracer = trace.Service("users.v1", logger)
	)

	return &UsersService{
		UnimplementedUsersServiceServer: pb.UnimplementedUsersServiceServer{},
		tracer:                          tracer,
		cases:                           New(uow),
	}
}

func (svc *UsersService) RegisterUser(
	ctx context.Context,
	request *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {
	user, err := svc.cases.RegisterUser(ctx, &RegisterUserInput{
		Name:     request.GetName(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, api.Error(err)
	}

	return &pb.RegisterUserResponse{User: MarshalUser(user)}, nil
}

func (svc *UsersService) RefreshToken(
	ctx context.Context,
	request *pb.RefreshTokenRequest,
) (*pb.RefreshTokenResponse, error) {
	user, err := svc.cases.RefreshToken(ctx, &RefreshTokenInput{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, api.Error(err)
	}

	return &pb.RefreshTokenResponse{User: MarshalUser(user)}, nil
}

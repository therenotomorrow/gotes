package v1

import (
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	pb "github.com/therenotomorrow/gotes/pkg/api/users/v1"
)

func MarshalUser(user *entities.User) *pb.User {
	return &pb.User{
		Id:    &typespb.ID{Value: user.ID.Value()},
		Name:  user.Name,
		Email: user.Email.Value(),
		Token: user.Token.Value(),
	}
}

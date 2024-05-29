package mock

import (
	"context"

	usv1 "github.com/zycgary/mxshop-go/internal/user/grpc/service/user/v1"
)

var _ usv1.UserRepository = (*userRepository)(nil)

type userRepository struct {
	users []*usv1.UserDO
}

func NewUserRepository() usv1.UserRepository {
	return &userRepository{}
}

func refreshUsers(ur *userRepository) {
	ur.users = make([]*usv1.UserDO, 0)
}

func (ur *userRepository) Index(ctx context.Context, page, pageSize int32) (*usv1.UserDOList, error) {
	refreshUsers(ur)

	ur.users = append(ur.users, &usv1.UserDO{
		Nickname: "test",
	})

	return &usv1.UserDOList{
		Total: int64(len(ur.users)),
		Data:  ur.users,
	}, nil
}

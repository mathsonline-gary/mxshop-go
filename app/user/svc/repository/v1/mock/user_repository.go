package mock

import (
	"context"
	urv1 "mxshop-go/app/user/svc/repository/v1"
)

type userRepository struct {
	users []*urv1.UserDO
}

func NewUserRepository() *userRepository {
	return &userRepository{}
}

func refreshUsers(ur *userRepository) {
	ur.users = make([]*urv1.UserDO, 0)
}

var _ urv1.UserRepository = (*userRepository)(nil)

func (ur *userRepository) Index(ctx context.Context, opts urv1.ListMeta) (*urv1.UserDOList, error) {
	refreshUsers(ur)

	ur.users = append(ur.users, &urv1.UserDO{
		Nickname: "test",
	})

	return &urv1.UserDOList{
		Total: int64(len(ur.users)),
		Data:  ur.users,
	}, nil
}

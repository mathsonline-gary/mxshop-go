package v1

import (
	"context"

	ucv1 "github.com/zycgary/mxshop-go/internal/user/grpc/controller/v1"
)

var _ ucv1.UserService = (*userService)(nil)

type UserDO struct {
	Nickname string `json:"name,omitempty"`
}

type UserDOList struct {
	Total int64     `json:"total,omitempty"`
	Data  []*UserDO `json:"data"`
}

type UserRepository interface {
	Index(ctx context.Context, page, pageSize int32) (*UserDOList, error)
}

type userService struct {
	userRepository UserRepository
}

func NewUserService(ur UserRepository) ucv1.UserService {
	return &userService{
		userRepository: ur,
	}
}

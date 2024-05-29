package v1

import (
	"context"

	pbv1 "github.com/zycgary/mxshop-go/api/user/v1"
)

var _ pbv1.UserServiceServer = (*userController)(nil)

type UserDTO struct {
	Nickname string `json:"name,omitempty"`
}

type UserDTOList struct {
	Total int64      `json:"total,omitempty"`
	Data  []*UserDTO `json:"data"`
}

type UserService interface {
	GetList(ctx context.Context, page, pageSize int32) (*UserDTOList, error)
}

type userController struct {
	pbv1.UnimplementedUserServiceServer
	us UserService
}

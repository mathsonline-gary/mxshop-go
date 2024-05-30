package v1

import (
	"context"

	ucv1 "github.com/zycgary/mxshop-go/internal/user/grpc/controller/v1"
	"github.com/zycgary/mxshop-go/pkg/log"
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
	logger         *log.Sugar
}

func NewUserService(ur UserRepository, logger log.Logger) ucv1.UserService {
	return &userService{
		userRepository: ur,
		logger:         log.NewSugar(logger),
	}
}

func (us *userService) GetList(ctx context.Context, page, pageSize int32) (*ucv1.UserListDTO, error) {
	ul, err := us.userRepository.Index(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var list ucv1.UserListDTO
	list.Total = ul.Total
	for _, v := range ul.Data {
		list.Data = append(list.Data, &ucv1.UserDTO{Nickname: v.Nickname})
	}

	return &list, nil
}

package v1

import (
	"context"
	udv1 "mxshop-go/app/user/svc/repository/v1"
)

type UserService interface {
	Index(context.Context, ListMeta) (*UserDTOList, error)
}

type userService struct {
	userRepository udv1.UserRepository
}

type ListMeta struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty" form:"page_size"`
}

type UserDTO struct {
	udv1.UserDO
}

type UserDTOList struct {
	Total int64      `json:"total,omitempty"`
	Data  []*UserDTO `json:"data"`
}

func NewUserService(ur udv1.UserRepository) *userService {
	return &userService{
		userRepository: ur,
	}
}

var _ UserService = (*userService)(nil)

func (us *userService) Index(ctx context.Context, opts ListMeta) (*UserDTOList, error) {
	ul, err := us.userRepository.Index(ctx, udv1.ListMeta{
		Page:     opts.Page,
		PageSize: opts.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var list UserDTOList
	list.Total = ul.Total
	for _, v := range ul.Data {
		list.Data = append(list.Data, &UserDTO{*v})
	}

	return &list, nil
}

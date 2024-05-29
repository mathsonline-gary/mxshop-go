package v1

import (
	"context"

	ucv1 "github.com/zycgary/mxshop-go/internal/user/grpc/controller/v1"
)

func (us *userService) GetList(ctx context.Context, page, pageSize int32) (*ucv1.UserDTOList, error) {
	ul, err := us.userRepository.Index(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var list ucv1.UserDTOList
	list.Total = ul.Total
	for _, v := range ul.Data {
		list.Data = append(list.Data, &ucv1.UserDTO{Nickname: v.Nickname})
	}

	return &list, nil
}

package data

import (
	"context"

	v1 "github.com/zycgary/mxshop-go/api/user/service/v1"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/logic"
	"github.com/zycgary/mxshop-go/pkg/log"
)

var _ logic.UserRepository = (*userRepository)(nil)

type userRepository struct {
	usc    v1.UserServiceClient
	logger *log.Sugar
}

func NewUserRepository(usc v1.UserServiceClient, logger log.Logger) logic.UserRepository {
	return &userRepository{
		usc:    usc,
		logger: log.NewSugar(logger),
	}
}

func (r *userRepository) GetList(ctx context.Context, page, pageSize int32) (*logic.UserList, error) {
	rsp, err := r.usc.GetUserList(ctx, &v1.GetUserListRequest{
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	})
	if err != nil {
		r.logger.Errorf("[Repository] [GetList]: %v", err)
		return nil, err
	}

	var ul logic.UserList
	ul.Total = rsp.Total
	for _, v := range rsp.Data {
		ul.Data = append(ul.Data, &logic.User{
			Nickname: v.Nickname,
		})
	}

	return &ul, nil
}

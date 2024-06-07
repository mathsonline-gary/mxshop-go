package service

import (
	"context"

	v1 "github.com/zycgary/mxshop-go/api/user/service/v1"
	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/logic"
	"github.com/zycgary/mxshop-go/pkg/log"
)

var _ v1.UserServiceServer = (*UserService)(nil)

type UserService struct {
	v1.UnimplementedUserServiceServer

	uc     *logic.UserUseCase
	logger *log.Sugar
}

func NewUserService(us *logic.UserUseCase, logger log.Logger) v1.UserServiceServer {
	return &UserService{
		uc:     us,
		logger: log.NewSugar(logger),
	}
}

func (s *UserService) GetUserList(ctx context.Context, req *v1.GetUserListRequest) (*v1.UserListResponse, error) {
	var page, pageSize int32 = 1, 10
	if req.Page > 0 {
		page = int32(req.Page)
	}
	if req.PageSize > 0 && req.PageSize <= 100 {
		pageSize = int32(req.PageSize)
	}

	list, err := s.uc.GetList(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var rsp v1.UserListResponse
	rsp.Total = list.Total
	for _, v := range list.Data {
		userInfo := v1.UserInfo{
			Nickname: v.Nickname,
		}
		rsp.Data = append(rsp.Data, &userInfo)
	}

	return &rsp, nil
}

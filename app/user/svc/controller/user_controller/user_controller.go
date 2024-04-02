package user_controller

import (
	"context"
	upb "mxshop-go/api/user/v1"
	usvc "mxshop-go/app/user/svc/service/v1"
)

type userController struct {
	upb.UnimplementedUserServiceServer
	us usvc.UserService
}

var _ upb.UserServiceServer = &userController{}

func (uc *userController) GetUserList(ctx context.Context, req *upb.GetUserListRequest) (*upb.UserListResponse, error) {
	opts := usvc.ListMeta{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}
	dtoList, err := uc.us.Index(ctx, opts)
	if err != nil {
		return nil, err
	}

	var rsp upb.UserListResponse
	rsp.Total = dtoList.Total
	for _, v := range dtoList.Data {
		userInfo := upb.UserInfo{
			Nickname: v.Nickname,
		}
		rsp.Data = append(rsp.Data, &userInfo)
	}

	return &rsp, nil
}

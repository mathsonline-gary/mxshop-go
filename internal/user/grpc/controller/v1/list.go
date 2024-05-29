package v1

import (
	"context"

	upb "github.com/zycgary/mxshop-go/api/user/v1"
)

func (uc *userController) GetUserList(ctx context.Context, req *upb.GetUserListRequest) (*upb.UserListResponse, error) {
	var page, pageSize int32 = 1, 10
	if req.Page > 0 {
		page = int32(req.Page)
	}
	if req.PageSize > 0 && req.PageSize <= 100 {
		pageSize = int32(req.PageSize)
	}

	dtoList, err := uc.us.GetList(ctx, page, pageSize)
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

package v1

import (
	"context"

	pbv1 "github.com/zycgary/mxshop-go/api/user/v1"
	"github.com/zycgary/mxshop-go/pkg/log"
)

var _ pbv1.UserServiceServer = (*userController)(nil)

type UserDTO struct {
	Nickname string `json:"name,omitempty"`
}

type UserListDTO struct {
	Total int64      `json:"total,omitempty"`
	Data  []*UserDTO `json:"data"`
}

type UserService interface {
	GetList(ctx context.Context, page, pageSize int32) (*UserListDTO, error)
}

type userController struct {
	pbv1.UnimplementedUserServiceServer

	us     UserService
	logger *log.Sugar
}

func NewUserController(us UserService, logger log.Logger) pbv1.UserServiceServer {
	return &userController{
		us:     us,
		logger: log.NewSugar(logger),
	}
}

func (uc *userController) GetUserList(ctx context.Context, req *pbv1.GetUserListRequest) (*pbv1.UserListResponse, error) {
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

	var rsp pbv1.UserListResponse
	rsp.Total = dtoList.Total
	for _, v := range dtoList.Data {
		userInfo := pbv1.UserInfo{
			Nickname: v.Nickname,
		}
		rsp.Data = append(rsp.Data, &userInfo)
	}

	return &rsp, nil
}

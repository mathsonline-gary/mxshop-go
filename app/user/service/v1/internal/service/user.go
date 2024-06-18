package service

import (
	"context"
	"crypto/sha512"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
	v1 "github.com/zycgary/mxshop-go/api/user/service/v1"
	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/logic"
	"github.com/zycgary/mxshop-go/pkg/log"
)

var _ v1.UserServiceServer = (*UserService)(nil)

var (
	passwordOptions = &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
)

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
			Id:       v.ID,
			Nickname: v.Nickname,
			Email:    v.Email,
		}
		rsp.Data = append(rsp.Data, &userInfo)
	}

	return &rsp, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *v1.EmailRequest) (*v1.UserInfoResponse, error) {
	user, err := s.uc.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &v1.UserInfoResponse{
		Data: &v1.UserInfo{
			Id:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
			Password: user.Password,
		},
	}, nil
}

func (s *UserService) CheckPassword(ctx context.Context, req *v1.CheckPasswordRequest) (*v1.CheckPasswordResponse, error) {
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], passwordOptions)

	return &v1.CheckPasswordResponse{
		Success: check,
	}, nil
}

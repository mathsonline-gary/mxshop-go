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
		ul.Data = append(ul.Data, &logic.SafeUser{
			ID:       v.Id,
			Nickname: v.Nickname,
			Email:    v.Email,
		})
	}

	return &ul, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*logic.User, error) {
	rsp, err := r.usc.GetUserByEmail(ctx, &v1.EmailRequest{
		Email: email,
	})
	if err != nil {
		r.logger.Errorf("[UserRepository] [GetUserByEmail]: %v", err)
		return nil, err
	}

	return &logic.User{
		SafeUser: logic.SafeUser{
			ID:       rsp.Data.Id,
			Email:    rsp.Data.Email,
			Nickname: rsp.Data.Nickname,
		},
		Password: rsp.Data.Password,
	}, nil
}

func (r *userRepository) CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error) {
	rsp, err := r.usc.CheckPassword(ctx, &v1.CheckPasswordRequest{
		Password:          password,
		EncryptedPassword: encryptedPassword,
	})
	if err != nil {
		r.logger.Errorf("[UserRepository] [CheckPassword]: %v", err)
		return false, err
	}

	return rsp.Success, nil
}

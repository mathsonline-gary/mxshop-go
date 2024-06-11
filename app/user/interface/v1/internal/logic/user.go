package logic

import (
	"context"

	"github.com/zycgary/mxshop-go/pkg/log"
)

type User struct {
	Nickname string `json:"name,omitempty"`
}

type UserList struct {
	Total int64   `json:"total,omitempty"`
	Data  []*User `json:"data"`
}

type UserRepository interface {
	GetList(ctx context.Context, page, pageSize int32) (*UserList, error)
}

type UserUseCase struct {
	ur     UserRepository
	logger *log.Sugar
}

func NewUserUseCase(ur UserRepository, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		ur:     ur,
		logger: log.NewSugar(logger),
	}
}

func (uc *UserUseCase) GetList(ctx context.Context, page, pageSize int32) (*UserList, error) {
	ul, err := uc.ur.GetList(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return ul, nil
}

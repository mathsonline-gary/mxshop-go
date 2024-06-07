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
	Index(ctx context.Context, page, pageSize int32) (*UserList, error)
}

type UserUseCase struct {
	userRepository UserRepository
	logger         *log.Sugar
}

func NewUserUseCase(ur UserRepository, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		userRepository: ur,
		logger:         log.NewSugar(logger),
	}
}

func (us *UserUseCase) GetList(ctx context.Context, page, pageSize int32) (*UserList, error) {
	ul, err := us.userRepository.Index(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return ul, nil
}

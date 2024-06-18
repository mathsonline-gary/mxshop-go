package logic

import (
	"context"

	"github.com/zycgary/mxshop-go/pkg/log"
)

type User struct {
	ID       uint64 `json:"id,omitempty"`
	Nickname string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserList struct {
	Total int64   `json:"total,omitempty"`
	Data  []*User `json:"data"`
}

type UserRepository interface {
	Index(ctx context.Context, page, pageSize int32) (*UserList, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
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
	return us.userRepository.Index(ctx, page, pageSize)
}

func (us *UserUseCase) GetByEmail(ctx context.Context, email string) (*User, error) {
	return us.userRepository.GetByEmail(ctx, email)
}

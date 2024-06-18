package logic

import (
	"context"

	"github.com/zycgary/mxshop-go/pkg/log"
)

type SafeUser struct {
	ID       uint64 `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}

type User struct {
	SafeUser

	Password string `json:"password,omitempty"`
}

type UserList struct {
	Total int64       `json:"total,omitempty"`
	Data  []*SafeUser `json:"data"`
}

type UserRepository interface {
	GetList(ctx context.Context, page, pageSize int32) (*UserList, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error)
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
	return uc.ur.GetList(ctx, page, pageSize)
}

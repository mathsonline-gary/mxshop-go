package mock

import (
	"context"

	usv1 "github.com/zycgary/mxshop-go/internal/user/grpc/service/user/v1"
	"gorm.io/gorm"
)

var _ usv1.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) usv1.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Index(ctx context.Context, page, pageSize int32) (*usv1.UserDOList, error) {
	return nil, nil
}

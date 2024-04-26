package mock

import (
	"context"
	urv1 "mxshop-go/app/user/svc/repository/v1"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() *userRepository {
	return &userRepository{}
}

var _ urv1.UserRepository = (*userRepository)(nil)

func (ur *userRepository) Index(ctx context.Context, opts urv1.ListMeta) (*urv1.UserDOList, error) {
	return nil, nil
}

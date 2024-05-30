package v1

import (
	"context"

	usv1 "github.com/zycgary/mxshop-go/internal/user/grpc/service/v1"
	"github.com/zycgary/mxshop-go/pkg/log"
	"gorm.io/gorm"
)

var _ usv1.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db     *gorm.DB
	logger *log.Sugar
}

func NewUserRepository(db *gorm.DB, logger log.Logger) usv1.UserRepository {
	return &userRepository{
		db:     db,
		logger: log.NewSugar(logger),
	}
}

func (ur *userRepository) Index(ctx context.Context, page, pageSize int32) (*usv1.UserDOList, error) {
	return nil, nil
}

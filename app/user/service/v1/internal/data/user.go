package data

import (
	"context"

	usv1 "github.com/zycgary/mxshop-go/app/user/service/v1/internal/logic"
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

func (ur *userRepository) Index(ctx context.Context, page, pageSize int32) (*usv1.UserList, error) {
	var total int64
	if err := ur.db.Model(&usv1.User{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return &usv1.UserList{
			Total: total,
			Data:  make([]*usv1.User, 0),
		}, nil
	}

	users := make([]*User, 0, pageSize)
	if err := ur.db.Scopes(paginate(page, pageSize)).Find(&users).Error; err != nil {
		return nil, err
	}
	ul := &usv1.UserList{
		Total: total,
		Data:  make([]*usv1.User, 0, pageSize),
	}
	for _, v := range users {
		userInfo := &usv1.User{
			ID:       uint64(v.ID),
			Nickname: v.Nickname,
			Email:    v.Email,
		}
		ul.Data = append(ul.Data, userInfo)
	}

	return ul, nil
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (*usv1.User, error) {
	var user usv1.User
	if err := ur.db.Where("email = ?", email).Find(&user).Limit(1).Error; err != nil {
		ur.logger.Errorf("[UserRepository] [GetByEmail]: %v", err)
		return nil, err
	}

	return &user, nil
}

func paginate(page, pageSize int32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

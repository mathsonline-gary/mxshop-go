package data

import (
	"errors"

	"github.com/zycgary/mxshop-go/user_svc/model"

	"gorm.io/gorm"
)

var _ UserRepo = (*userRepo)(nil)

type UserRepo interface {
	ListUser(page int32, pageSize int32) (total int64, users []*model.User, err error)
	GetUserByID(id int32) (*model.User, error)
	GetUserByMobile(mobile string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (u userRepo) ListUser(page int32, pageSize int32) (int64, []*model.User, error) {
	var total int64
	if err := u.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return 0, make([]*model.User, 0), nil
	}

	users := make([]*model.User, 0, pageSize)

	if err := u.db.Scopes(paginate(page, pageSize)).Find(&users).Error; err != nil {
		return 0, nil, err
	}

	return total, users, nil
}

func (u userRepo) GetUserByID(id int32) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) GetUserByMobile(mobile string) (*model.User, error) {
	var user model.User
	result := u.db.Where("mobile = ?", mobile).Limit(1).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (u userRepo) CreateUser(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) UpdateUser(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func paginate(page, pageSize int32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

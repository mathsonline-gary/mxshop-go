package data

import (
	"mxshop-go/user_svc/model"

	"gorm.io/gorm"
)

var _ UserRepo = (*userRepo)(nil)

type UserRepo interface {
	ListUser(page uint32, pageSize uint32) (total int64, users []*model.User, err error)
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

func (u userRepo) ListUser(page uint32, pageSize uint32) (int64, []*model.User, error) {
	var total int64
	users := make([]*model.User, 0, pageSize)

	if err := u.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return 0, users, err
	}

	if err := u.db.Scopes(Paginate(int(page), int(pageSize))).Find(&users).Error; err != nil {
		return 0, users, err
	}

	return total, users, nil
}

func (u userRepo) GetUserByID(id int32) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) GetUserByMobile(mobile string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) CreateUser(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) UpdateUser(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

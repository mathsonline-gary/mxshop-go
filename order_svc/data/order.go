package data

import (
	"context"
	"errors"

	"github.com/zycgary/mxshop-go/order_svc/model"
	"gorm.io/gorm"
)

var _ OrderRepo = (*orderRepo)(nil)

type OrderRepo interface {
	ListCartItems(context.Context, int32) ([]*model.CartItem, error)
	GetCartItemByProductID(context.Context, int32, int32) (*model.CartItem, error)
	GetCartItemByID(context.Context, int32) (*model.CartItem, error)
	UpsertCartItem(context.Context, *model.CartItem) error
	DeleteCartItem(context.Context, int32) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) ListCartItems(ctx context.Context, userId int32) ([]*model.CartItem, error) {
	var items []*model.CartItem
	if err := r.db.Where("user_id = ?", userId).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *orderRepo) GetCartItemByProductID(ctx context.Context, userId, productId int32) (*model.CartItem, error) {
	var item model.CartItem
	if err := r.db.Where("user_id = ? AND product_id = ?", userId, productId).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &item, nil
}

func (r *orderRepo) GetCartItemByID(ctx context.Context, id int32) (*model.CartItem, error) {
	var item model.CartItem
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &item, nil
}

func (r *orderRepo) UpsertCartItem(ctx context.Context, item *model.CartItem) error {
	return r.db.Save(item).Error
}

func (r *orderRepo) DeleteCartItem(ctx context.Context, id int32) error {
	return r.db.Delete(&model.CartItem{}, id).Error
}

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
	CountOrders(context.Context, int32) (int64, error)
	ListOrders(context.Context, int32, int32, int32) ([]*model.Order, error)
	GetOrderByID(context.Context, int32) (*model.Order, error)
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

func (r *orderRepo) CountOrders(ctx context.Context, userId int32) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Order{}).Where(&model.Order{UserID: userId}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *orderRepo) ListOrders(ctx context.Context, userId, page, pageSize int32) ([]*model.Order, error) {
	orders := make([]*model.Order, 0, pageSize)
	if err := r.db.Scopes(paginate(page, pageSize)).Where(&model.Order{UserID: userId}).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepo) GetOrderByID(ctx context.Context, id int32) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("Items").First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &order, nil
}

func paginate(page, pageSize int32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

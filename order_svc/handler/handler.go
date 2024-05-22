package handler

import (
	"github.com/zycgary/mxshop-go/order_svc/model"
	"github.com/zycgary/mxshop-go/order_svc/proto"
)

var _ proto.OrderServiceServer = (*OrderService)(nil)

type Option func(*OrderService)

type OrderService struct {
	proto.UnimplementedOrderServiceServer

	repo model.OrderRepo
}

func NewOrderService(opts ...Option) *OrderService {
	s := &OrderService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithRepo(repo model.OrderRepo) Option {
	return func(o *OrderService) {
		o.repo = repo
	}
}

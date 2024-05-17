package handler

import (
	"github.com/zycgary/mxshop-go/order_svc/model"
	"github.com/zycgary/mxshop-go/order_svc/proto"
)

var _ proto.OrderServiceServer = (*orderServiceServer)(nil)

type Option func(*orderServiceServer)

type orderServiceServer struct {
	proto.UnimplementedOrderServiceServer

	repo model.OrderRepo
}

func NewOrderServiceServer(opts ...Option) proto.OrderServiceServer {
	s := &orderServiceServer{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithRepo(repo model.OrderRepo) Option {
	return func(o *orderServiceServer) {
		o.repo = repo
	}
}

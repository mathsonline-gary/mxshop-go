package handler

import (
	"context"

	"github.com/zycgary/mxshop-go/order_svc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *OrderServiceServer) CreateOrder(ctx context.Context, request *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, request *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, request *proto.GetOrderRequest) (*proto.GetOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, request *proto.UpdateOrderStatusRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

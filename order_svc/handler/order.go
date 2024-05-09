package handler

import (
	"context"

	"github.com/zycgary/mxshop-go/order_svc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *OrderServiceServer) CreateOrder(ctx context.Context, request *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, request *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	total, err := s.orderRepo.CountOrders(ctx, request.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	orders, err := s.orderRepo.ListOrders(ctx, request.UserId, request.Page, request.PageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	rsp := &proto.ListOrdersResponse{
		Total: total,
		Data:  make([]*proto.Order, 0, len(orders)),
	}
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.Order{
			Id:              order.ID,
			UserId:          order.UserID,
			SerialNumber:    order.SerialNumber,
			PaymentMethod:   order.PaymentMethod,
			Status:          order.Status,
			Amount:          order.Amount,
			ShippingAddress: order.ShippingAddress,
			ShippingName:    order.ShippingName,
			ShippingMobile:  order.ShippingMobile,
			Note:            order.Note,
			TradingNumber:   order.TradingNumber,
			PaidAt:          timestamppb.New(order.PaidAt),
			Items:           nil,
		})
	}

	return rsp, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, request *proto.GetOrderRequest) (*proto.GetOrderResponse, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	if order == nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	rsp := &proto.GetOrderResponse{
		Data: &proto.Order{
			Id:              order.ID,
			UserId:          order.UserID,
			SerialNumber:    order.SerialNumber,
			PaymentMethod:   order.PaymentMethod,
			Status:          order.Status,
			Amount:          order.Amount,
			ShippingAddress: order.ShippingAddress,
			ShippingName:    order.ShippingName,
			ShippingMobile:  order.ShippingMobile,
			Note:            order.Note,
			TradingNumber:   order.TradingNumber,
			PaidAt:          timestamppb.New(order.PaidAt),
			Items:           make([]*proto.OrderItem, 0, len(order.Items)),
		},
	}
	for _, item := range order.Items {
		rsp.Data.Items = append(rsp.Data.Items, &proto.OrderItem{
			Id:                item.ID,
			OrderId:           item.OrderID,
			ProductId:         item.ProductID,
			ProductName:       item.ProductName,
			ProductImage:      item.ProductImage,
			ProductUnitPrice:  item.ProductUnitPrice,
			Quantity:          item.Quantity,
			ProductTotalPrice: item.ProductTotalPrice,
		})
	}

	return rsp, nil
}

func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, request *proto.UpdateOrderStatusRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

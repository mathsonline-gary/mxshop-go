package handler

import (
	"context"

	"github.com/zycgary/mxshop-go/order_svc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *OrderService) CreateOrder(ctx context.Context, request *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	// get selected products from cart
	items, err := s.repo.ListCartItems(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	pids := make([]int32, 0, len(items))
	for _, item := range items {
		if item.Selected {
			pids = append(pids, item.ProductID)
		}
	}
	if len(pids) == 0 {
		return nil, status.Error(codes.InvalidArgument, proto.ErrorNoSelectedProduct)
	}

	// TODO: get product details via product service

	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

// ListOrders Retrieves a paginated list of orders for a user.
func (s *OrderService) ListOrders(ctx context.Context, request *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	total, err := s.repo.CountOrders(ctx, request.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, proto.ErrorInternal)
	}

	orders, err := s.repo.ListOrders(ctx, request.UserId, request.Page, request.PageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, proto.ErrorInternal)
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

func (s *OrderService) GetOrder(ctx context.Context, request *proto.GetOrderRequest) (*proto.GetOrderResponse, error) {
	order, err := s.repo.GetOrderByID(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, proto.ErrorInternal)
	}
	if order == nil {
		return nil, status.Error(codes.NotFound, proto.ErrorOrderNotFound)
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

// UpdateOrderStatus Updates the status of an order by its serial number.
func (s *OrderService) UpdateOrderStatus(ctx context.Context, request *proto.UpdateOrderStatusRequest) (*emptypb.Empty, error) {
	err := s.repo.UpdateOrderStatus(ctx, request.SerialNumber, request.Status)
	if err != nil {
		if err.Error() == proto.ErrorOrderNotFound {
			return nil, status.Error(codes.NotFound, proto.ErrorOrderNotFound)
		}
		return nil, status.Error(codes.Internal, proto.ErrorInternal)
	}

	return &emptypb.Empty{}, nil
}

package handler

import (
	"context"

	"github.com/zycgary/mxshop-go/order_svc/model"
	"github.com/zycgary/mxshop-go/order_svc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListCartItems retrieves the list of shopping cart items for a user.
func (s *OrderServiceServer) ListCartItems(ctx context.Context, request *proto.ListCartItemsRequest) (*proto.ListCartItemsResponse, error) {
	uid := request.UserId
	if uid <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	items, err := s.orderRepo.ListCartItems(ctx, uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list cart items")
	}

	rsp := &proto.ListCartItemsResponse{
		Total: int32(len(items)),
		Data:  make([]*proto.CartItem, 0, len(items)),
	}
	for _, item := range items {
		rsp.Data = append(rsp.Data, &proto.CartItem{
			Id:        item.ID,
			UserId:    item.UserID,
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Selected:  item.Selected,
		})
	}

	return rsp, nil
}

// AddCartItem adds given product(s) to the user's shopping cart.
// If the product already exists in the cart, it will increase the quantity.
func (s *OrderServiceServer) AddCartItem(ctx context.Context, request *proto.AddCartItemRequest) (*proto.AddCartItemResponse, error) {
	item, err := s.orderRepo.GetCartItemByProductID(ctx, request.UserId, request.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add cart item.")
	}
	if item != nil {
		// product already exists in the cart, increase the quantity
		item.Quantity += request.Quantity
	} else {
		// product does not exist in the cart, add it
		item = &model.CartItem{
			UserID:    request.UserId,
			ProductID: request.ProductId,
			Quantity:  request.Quantity,
			Selected:  true,
		}
	}

	if err := s.orderRepo.UpsertCartItem(ctx, item); err != nil {
		return nil, err
	}

	return &proto.AddCartItemResponse{
		Data: &proto.CartItem{
			Id:        item.ID,
			UserId:    item.UserID,
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Selected:  item.Selected,
		},
	}, nil
}

// UpdateCartItem updates a user's shopping cart item.
// It will update the quantity and selected status of the cart item
func (s *OrderServiceServer) UpdateCartItem(ctx context.Context, request *proto.UpdateCartItemRequest) (*emptypb.Empty, error) {
	// Validate the request
	if request.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid cart item ID")
	}
	if request.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid quantity")
	}

	// Get the cart item
	item, err := s.orderRepo.GetCartItemByID(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	if item == nil {
		return nil, status.Errorf(codes.NotFound, "cart item not found")
	}

	// Update the quantity
	item.Quantity = request.Quantity
	item.Selected = request.Selected
	if err := s.orderRepo.UpsertCartItem(ctx, item); err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}

// DeleteCartItem deletes a cart item from the user's shopping cart.
func (s *OrderServiceServer) DeleteCartItem(ctx context.Context, request *proto.DeleteCartItemRequest) (*emptypb.Empty, error) {
	if err := s.orderRepo.DeleteCartItem(ctx, request.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}

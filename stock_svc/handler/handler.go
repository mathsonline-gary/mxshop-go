package handler

import (
	"context"

	"mxshop-go/product_svc/global"
	"mxshop-go/stock_svc/model"
	"mxshop-go/stock_svc/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StockServiceServer struct {
	proto.UnimplementedStockServiceServer
}

var _ proto.StockServiceServer = (*StockServiceServer)(nil)

func (s StockServiceServer) UpsertStock(ctx context.Context, request *proto.UpsertStockRequest) (*emptypb.Empty, error) {
	if request.ProductId <= 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid product id")
	}
	if request.Quantity < 0 {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid quantity")
	}

	var stock model.Stock

	if err := global.DB.Limit(1).Where("product_id = ?", request.ProductId).Find(&stock).Error; err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, err.Error())
	}
	stock.ProductID = request.ProductId
	stock.Quantity = request.Quantity
	global.DB.Save(&stock)

	return &emptypb.Empty{}, nil
}

func (s StockServiceServer) GetStock(ctx context.Context, request *proto.GetStockRequest) (*proto.GetStockResponse, error) {
	if request.ProductId <= 0 {
		return nil, status.Errorf(codes.NotFound, "stock not found")
	}

	var stock model.Stock

	result := global.DB.Limit(1).Where("product_id = ?", request.ProductId).Find(&stock)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "stock not found")
	}

	return &proto.GetStockResponse{
		ProductId: stock.ProductID,
		Quantity:  stock.Quantity,
	}, nil
}

func (s StockServiceServer) Withhold(ctx context.Context, request *proto.WithholdRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s StockServiceServer) Return(ctx context.Context, request *proto.ReturnRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

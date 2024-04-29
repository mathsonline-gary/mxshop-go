package handler

import (
	"context"

	"mxshop-go/stock_svc/global"
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
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID")
	}
	if request.Quantity < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid quantity")
	}

	var stock model.Stock

	if err := global.DB.Limit(1).Where(&model.Stock{ProductID: request.ProductId}).Find(&stock).Error; err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
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

func (s StockServiceServer) WithholdStock(ctx context.Context, request *proto.WithholdStockRequest) (*emptypb.Empty, error) {
	// validate request
	for _, data := range request.Data {
		if data.ProductId <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "product not found")
		}
		if data.Quantity < 0 {
			return nil, status.Errorf(codes.InvalidArgument, "invalid quantity")
		}
	}

	// withhold stock
	tx := global.DB.Begin()

	for _, data := range request.Data {
		var stock model.Stock

		// get stock
		result := global.DB.Limit(1).Where(&model.Stock{ProductID: data.ProductId}).Find(&stock)
		if result.Error != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "product not found")
		}
		// check stock quantity
		if stock.Quantity < data.Quantity {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "insufficient stock")
		}

		// update stock
		stock.ProductID = data.ProductId
		stock.Quantity -= data.Quantity
		tx.Save(&stock)
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}

// ReturnStock returns the withheld stocks. It may be called when:
// 1. order is timeout
// 2. order is cancelled
// 3. order is failed to be created
func (s StockServiceServer) ReturnStock(ctx context.Context, request *proto.ReturnStockRequest) (*emptypb.Empty, error) {
	// validate request
	for _, data := range request.Data {
		if data.ProductId <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "product not found")
		}
		if data.Quantity < 0 {
			return nil, status.Errorf(codes.InvalidArgument, "invalid quantity")
		}
	}

	// return stock
	tx := global.DB.Begin()

	for _, data := range request.Data {
		var stock model.Stock

		// get stock
		result := global.DB.Limit(1).Where(&model.Stock{ProductID: data.ProductId}).Find(&stock)
		if result.Error != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "product not found")
		}

		// update stock
		stock.ProductID = data.ProductId
		stock.Quantity += data.Quantity
		tx.Save(&stock)
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}

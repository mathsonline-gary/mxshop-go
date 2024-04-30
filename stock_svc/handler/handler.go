package handler

import (
	"context"
	"slices"

	"mxshop-go/stock_svc/global"
	"mxshop-go/stock_svc/model"
	"mxshop-go/stock_svc/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StockServiceServer struct {
	proto.UnimplementedStockServiceServer
}

var _ proto.StockServiceServer = (*StockServiceServer)(nil)

func (s StockServiceServer) UpsertStock(_ context.Context, request *proto.UpsertStockRequest) (*emptypb.Empty, error) {
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

func (s StockServiceServer) GetStock(_ context.Context, request *proto.GetStockRequest) (*proto.GetStockResponse, error) {
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

func (s StockServiceServer) WithholdStock(_ context.Context, request *proto.WithholdStockRequest) (*emptypb.Empty, error) {
	// validate request
	for _, data := range request.Data {
		if data.ProductId <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "product not found")
		}
		if data.Quantity < 0 {
			return nil, status.Errorf(codes.InvalidArgument, "invalid quantity")
		}
	}

	// sort request data by product ID to acquire pessimistic locks in a consistent order, which can help prevent deadlocks
	if len(request.Data) > 1 {
		slices.SortFunc(request.Data, func(a, b *proto.StockInfo) int {
			if a.ProductId < b.ProductId {
				return -1
			}
			return 1
		})
	}

	// withhold stock
	tx := global.DB.Begin()

	for _, data := range request.Data {
		var stock model.Stock
		// Optimistic Locking: retry if failed to update stock caused by lock
		// This approach is can can be effective in scenarios where conflicts are rare, i.e., it's uncommon for multiple transactions to try to update the same stock item at the same time.
		for {
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

			result = tx.Model(&model.Stock{}).Where("product_id = ? and version = ?", data.ProductId, stock.Version).Select("Quantity", "Version").Updates(model.Stock{Quantity: stock.Quantity, Version: stock.Version + 1})
			if result.Error != nil {
				tx.Rollback()
				return nil, status.Errorf(codes.Internal, result.Error.Error())
			}
			if result.RowsAffected == 0 {
				zap.S().Info("failed to update stock caused by lock, retrying...")
				continue
			}
			break
		}

		// Pessimistic Locking: lock the stock item before updating it
		// This approach is effective in scenarios where conflicts are common, i.e., it's common for multiple transactions to try to update the same stock item at the same time.
		// But it can potentially decrease throughput and increase risk of deadlocks.
		// To mitigate the risk of deadlocks,  one common practice is to always acquire locks in a consistent order. For example, we could sort the product IDs before starting the transaction, and then always update the products in the order of their IDs. This ensures that all transactions acquire locks in the same order, which can help prevent deadlocks.
		/*
			result := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Where(&model.Stock{ProductID: data.ProductId}).Find(&stock)
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
		*/
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}

// ReturnStock returns the withheld stocks. It may be called when:
// 1. order is timeout
// 2. order is cancelled
// 3. order is failed to be created
func (s StockServiceServer) ReturnStock(_ context.Context, request *proto.ReturnStockRequest) (*emptypb.Empty, error) {
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

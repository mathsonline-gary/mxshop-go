package tests

import (
	"context"
	"sync"
	"testing"

	"mxshop-go/stock_svc/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUpsertStock(t *testing.T) {
	type args struct {
		ctx context.Context
		req *proto.UpsertStockRequest
	}

	type expected struct {
		rsp *emptypb.Empty
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "create stock for a new product",
			args: args{
				ctx: context.Background(),
				req: &proto.UpsertStockRequest{
					ProductId: 99999,
					Quantity:  10,
				},
			},
			expected: expected{
				rsp: &emptypb.Empty{},
				err: nil,
			},
		},
		{
			name: "update an existing stock",
			args: args{
				ctx: context.Background(),
				req: &proto.UpsertStockRequest{
					ProductId: 421,
					Quantity:  10,
				},
			},
			expected: expected{
				rsp: &emptypb.Empty{},
				err: nil,
			},
		},
		{
			name: "invalid product id",
			args: args{
				ctx: context.Background(),
				req: &proto.UpsertStockRequest{
					ProductId: 0,
					Quantity:  10,
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.InvalidArgument, "invalid product ID"),
			},
		},
		{
			name: "invalid quantity",
			args: args{
				ctx: context.Background(),
				req: &proto.UpsertStockRequest{
					ProductId: 420,
					Quantity:  -10,
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.InvalidArgument, "invalid quantity"),
			},
		},
	}

	for _, tt := range tests {
		got, err := stockClient.UpsertStock(tt.args.ctx, tt.args.req)

		if tt.expected.err != nil {
			if err == nil {
				t.Errorf(`test "%s" failed, expect error: %v, got: <nil>`, tt.name, tt.expected.err)
				continue
			}
			if status.Code(err) != status.Code(tt.expected.err) || err.Error() != tt.expected.err.Error() {
				t.Errorf(`test "%s" failed, expected error: %v, got: %v`, tt.name, tt.expected.err, err)
				continue
			}
		}

		if tt.expected.rsp != nil {
			if got == nil {
				t.Errorf(`test "%s" failed, expected response: &emptypb.Empty{}, got: <nil>`, tt.name)
				continue
			}
		}
	}
}

func TestGetStock(t *testing.T) {
	type args struct {
		ctx context.Context
		req *proto.GetStockRequest
	}

	type expected struct {
		rsp *proto.GetStockResponse
		err error
	}

	test := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "get stock of a product",
			args: args{
				ctx: context.Background(),
				req: &proto.GetStockRequest{
					ProductId: 421,
				},
			},
			expected: expected{
				rsp: &proto.GetStockResponse{
					ProductId: 421,
				},
				err: nil,
			},
		},
		{
			name: "invalid product ID",
			args: args{
				ctx: context.Background(),
				req: &proto.GetStockRequest{
					ProductId: 0,
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.NotFound, "stock not found"),
			},
		},
		{
			name: "non-existent product",
			args: args{

				ctx: context.Background(),
				req: &proto.GetStockRequest{
					ProductId: 9999,
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.NotFound, "stock not found"),
			},
		},
	}

	for _, tt := range test {
		got, err := stockClient.GetStock(tt.args.ctx, tt.args.req)

		if (tt.expected.rsp == nil) != (got == nil) || (tt.expected.err == nil) != (err == nil) {
			t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
			continue
		}

		if tt.expected.err != nil {
			if status.Code(err) != status.Code(tt.expected.err) || err.Error() != tt.expected.err.Error() {
				t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
				continue
			}
		}

		if tt.expected.rsp != nil {
			if got == nil || got.ProductId != tt.expected.rsp.ProductId {
				t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
				continue
			}
		}
	}
}

func TestWithholdStock(t *testing.T) {
	type args struct {
		ctx context.Context
		req *proto.WithholdStockRequest
	}

	type expected struct {
		rsp *emptypb.Empty
		err error
	}

	test := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "withhold stock of single product",
			args: args{
				ctx: context.Background(),
				req: &proto.WithholdStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 840,
							Quantity:  5,
						},
					},
				},
			},
			expected: expected{
				rsp: &emptypb.Empty{},
				err: nil,
			},
		},
		{
			name: "withhold stock of multiple products",
			args: args{
				ctx: context.Background(),
				req: &proto.WithholdStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 840,
							Quantity:  5,
						},
						{
							ProductId: 839,
							Quantity:  10,
						},
					},
				},
			},
			expected: expected{
				rsp: &emptypb.Empty{},
				err: nil,
			},
		},
		{
			name: "invalid product id",
			args: args{
				ctx: context.Background(),
				req: &proto.WithholdStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 0,
							Quantity:  5,
						},
					},
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.InvalidArgument, "product not found"),
			},
		},
		{
			name: "invalid quantity",
			args: args{
				ctx: context.Background(),
				req: &proto.WithholdStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 840,
							Quantity:  -50,
						},
					},
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.InvalidArgument, "invalid quantity"),
			},
		},
		{
			name: "non-existent product",
			args: args{
				ctx: context.Background(),
				req: &proto.WithholdStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 9999,
							Quantity:  50,
						},
					},
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.InvalidArgument, "product not found"),
			},
		},
		{
			name: "not enough stock",
			args: args{
				ctx: context.Background(),
				req: &proto.WithholdStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 840,
							Quantity:  99999,
						},
					},
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.ResourceExhausted, "insufficient stock"),
			},
		},
	}

	for _, tt := range test {
		got, err := stockClient.WithholdStock(tt.args.ctx, tt.args.req)

		if (tt.expected.rsp == nil) != (got == nil) || (tt.expected.err == nil) != (err == nil) {
			t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
			continue
		}

		if tt.expected.err != nil {
			if status.Code(err) != status.Code(tt.expected.err) || err.Error() != tt.expected.err.Error() {
				t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
				continue
			}
		}

		if tt.expected.rsp != nil {
			if got == nil {
				t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
				continue
			}
		}
	}
}

func TestWithholdStockConcurrency(t *testing.T) {
	l := 20
	var wg sync.WaitGroup
	wg.Add(l)
	for i := 0; i < l; i++ {
		go func() {
			defer wg.Done()
			_, err := stockClient.WithholdStock(context.Background(), &proto.WithholdStockRequest{
				Data: []*proto.StockInfo{
					{
						ProductId: 840,
						Quantity:  1,
					},
					{
						ProductId: 839,
						Quantity:  1,
					},
				},
			})
			if err != nil {
				t.Errorf("error: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestReturnStock(t *testing.T) {
	type args struct {
		ctx context.Context
		req *proto.ReturnStockRequest
	}

	type expected struct {
		rsp *emptypb.Empty
		err error
	}

	test := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "return stock of single product",
			args: args{
				ctx: context.Background(),
				req: &proto.ReturnStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 840,
							Quantity:  5,
						},
					},
				},
			},
			expected: expected{
				rsp: &emptypb.Empty{},
				err: nil,
			},
		},
		{
			name: "return stock of multiple products",
			args: args{
				ctx: context.Background(),
				req: &proto.ReturnStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 840,
							Quantity:  5,
						},
						{
							ProductId: 839,
							Quantity:  10,
						},
					},
				},
			},
			expected: expected{
				rsp: &emptypb.Empty{},
				err: nil,
			},
		},
		{
			name: "invalid product ID",
			args: args{
				ctx: context.Background(),
				req: &proto.ReturnStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 0,
							Quantity:  5,
						},
					},
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.InvalidArgument, "product not found"),
			},
		},
		{
			name: "invalid quantity",
			args: args{
				ctx: context.Background(),
				req: &proto.ReturnStockRequest{
					Data: []*proto.StockInfo{
						{
							ProductId: 840,
							Quantity:  -5,
						},
					},
				},
			},
			expected: expected{
				rsp: nil,
				err: status.Errorf(codes.InvalidArgument, "invalid quantity"),
			},
		},
	}

	for _, tt := range test {
		got, err := stockClient.ReturnStock(tt.args.ctx, tt.args.req)

		if (tt.expected.rsp == nil) != (got == nil) || (tt.expected.err == nil) != (err == nil) {
			t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
			continue
		}

		if tt.expected.err != nil {
			if status.Code(err) != status.Code(tt.expected.err) || err.Error() != tt.expected.err.Error() {
				t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
				continue
			}
		}

		if tt.expected.rsp != nil {
			if got == nil {
				t.Errorf(`test "%s" failed, expected: (%+v, %v), got: (%+v, %v)`, tt.name, tt.expected.rsp, tt.expected.err, got, err)
				continue
			}
		}
	}
}

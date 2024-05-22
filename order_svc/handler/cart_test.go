package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zycgary/mxshop-go/order_svc/data/mock"
	"github.com/zycgary/mxshop-go/order_svc/model"
	"github.com/zycgary/mxshop-go/order_svc/proto"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestOrderServiceServer_ListCartItems(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *proto.ListCartItemsRequest
	}

	type asserts func(*testing.T, *proto.ListCartItemsResponse, error)

	tests := []struct {
		name    string
		args    args
		expects func(repo *mock.MockOrderRepo)
		asserts asserts
	}{
		{
			name: "Invalid user ID",
			args: args{
				ctx: context.Background(),
				request: &proto.ListCartItemsRequest{
					UserId: 0,
				},
			},
			expects: func(m *mock.MockOrderRepo) {},
			asserts: func(t *testing.T, response *proto.ListCartItemsResponse, err error) {
				assert.Nil(t, response)
				assert.Error(t, err)
				assert.ErrorIs(t, err, status.Errorf(codes.InvalidArgument, "invalid user ID"))
			},
		},
		{
			name: "Internal error",
			args: args{
				ctx: context.Background(),
				request: &proto.ListCartItemsRequest{
					UserId: 1,
				},
			},
			expects: func(m *mock.MockOrderRepo) {
				m.EXPECT().ListCartItems(context.Background(), int32(1)).Return(nil, errors.New("internal error"))
			},

			asserts: func(t *testing.T, response *proto.ListCartItemsResponse, err error) {
				assert.Nil(t, response)
				assert.Error(t, err)
				assert.ErrorIs(t, err, status.Errorf(codes.Internal, "failed to list cart items"))
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				request: &proto.ListCartItemsRequest{
					UserId: 1,
				},
			},
			expects: func(m *mock.MockOrderRepo) {
				orders := make([]*model.CartItem, 0, 10)
				for i := 0; i < 10; i++ {
					orders = append(orders, &model.CartItem{})
				}
				m.EXPECT().ListCartItems(context.Background(), int32(1)).Return(orders, nil)
			},

			asserts: func(t *testing.T, response *proto.ListCartItemsResponse, err error) {
				assert.NotNil(t, response)
				assert.Equal(t, int32(10), response.Total)
				assert.Len(t, response.Data, 10)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock.NewMockOrderRepo(ctrl)
			server := NewOrderService(WithRepo(m))

			tt.expects(m)
			got, err := server.ListCartItems(tt.args.ctx, tt.args.request)
			tt.asserts(t, got, err)
		})
	}
}

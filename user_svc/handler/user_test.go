package handler

import (
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-go/user_svc/data/mock"
	"mxshop-go/user_svc/model"
	"mxshop-go/user_svc/proto"
	"testing"

	userproto "mxshop-go/user_svc/proto"
)

func TestUserServiceServer_GetUserList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock.NewMockUserRepo(ctrl)
	server := NewUserServiceServer(mockUserRepo)

	type args struct {
		ctx     context.Context
		request *userproto.GetUserListRequest
	}

	type expect struct {
		rsp *userproto.UserListResponse
		err error
	}

	tests := []struct {
		name   string
		args   args
		mock   func()
		expect expect
	}{
		{
			name: "Internal error",
			args: args{
				ctx: context.Background(),
				request: &userproto.GetUserListRequest{
					Page:     1,
					PageSize: 10,
				},
			},
			mock: func() {
				var total int64 = 0
				mockUserRepo.EXPECT().ListUser(gomock.Any(), gomock.Any()).Return(total, nil, errors.New("internal error"))
			},
			expect: expect{
				rsp: nil,
				err: status.Errorf(codes.Internal, "get user list"),
			},
		},
		{
			name: "Empty user list",
			args: args{
				ctx: context.Background(),
				request: &userproto.GetUserListRequest{
					Page:     1,
					PageSize: 10,
				},
			},
			mock: func() {
				var total int64 = 0
				users := make([]*model.User, 0)
				mockUserRepo.EXPECT().ListUser(gomock.Any(), gomock.Any()).Return(total, users, nil)
			},
			expect: expect{
				rsp: &userproto.UserListResponse{
					Total: 0,
					Data:  []*proto.UserInfo{},
				},
				err: nil,
			},
		},
		{
			name: "User list",
			args: args{
				ctx: context.Background(),
				request: &userproto.GetUserListRequest{
					Page:     1,
					PageSize: 10,
				},
			},
			mock: func() {
				var total int64 = 20
				users := make([]*model.User, 0, 10)
				for i := 0; i < 10; i++ {
					users = append(users, &model.User{})
				}
				mockUserRepo.EXPECT().ListUser(gomock.Any(), gomock.Any()).Return(total, users, nil)
			},
			expect: expect{
				rsp: &userproto.UserListResponse{
					Total: 20,
					Data:  make([]*userproto.UserInfo, 10),
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := server.GetUserList(tt.args.ctx, tt.args.request)
			if !errors.Is(err, tt.expect.err) {
				t.Errorf("GetUserList() got = (%v, %v), expect (%v, %v)", got, err, tt.expect.rsp, tt.expect.err)
				return
			}
			if err == nil {
				if got.Total != tt.expect.rsp.Total || len(got.Data) != len(tt.expect.rsp.Data) {
					t.Errorf("GetUserList() got = (%v, %v), expect (%v, %v)", got, err, tt.expect.rsp, tt.expect.err)
					return
				}
			}
		})
	}
}

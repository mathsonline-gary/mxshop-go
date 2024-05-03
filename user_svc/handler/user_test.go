package handler

import (
	"context"
	"errors"
	"testing"

	"mxshop-go/user_svc/data/mock"
	"mxshop-go/user_svc/model"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userproto "mxshop-go/user_svc/proto"
)

func TestUserServiceServer_GetUserList(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *userproto.GetUserListRequest
	}

	type asserts func(*testing.T, *userproto.UserListResponse, error)

	tests := []struct {
		name    string
		args    args
		expects func(*mock.MockUserRepo)
		asserts asserts
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
			expects: func(mockUserRepo *mock.MockUserRepo) {
				var total int64 = 0
				mockUserRepo.EXPECT().ListUser(gomock.Any(), gomock.Any()).Return(total, nil, errors.New("internal error"))
			},
			asserts: func(t *testing.T, response *userproto.UserListResponse, err error) {
				assert.Nil(t, response)
				assert.Error(t, err)
				assert.ErrorIs(t, err, status.Errorf(codes.Internal, "get user list"))
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
			expects: func(mockUserRepo *mock.MockUserRepo) {
				var total int64 = 0
				users := make([]*model.User, 0)
				mockUserRepo.EXPECT().ListUser(gomock.Any(), gomock.Any()).Return(total, users, nil)
			},
			asserts: func(t *testing.T, response *userproto.UserListResponse, err error) {
				assert.NotNil(t, response)
				assert.Equal(t, int64(0), response.Total)
				assert.Empty(t, response.Data)
				assert.NoError(t, err)
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
			expects: func(mockUserRepo *mock.MockUserRepo) {
				var total int64 = 20
				users := make([]*model.User, 0, 10)
				for i := 0; i < 10; i++ {
					users = append(users, &model.User{})
				}
				mockUserRepo.EXPECT().ListUser(gomock.Any(), gomock.Any()).Return(total, users, nil)
			},
			asserts: func(t *testing.T, response *userproto.UserListResponse, err error) {
				assert.NotNil(t, response)
				assert.Equal(t, int64(20), response.Total)
				assert.Len(t, response.Data, 10)
				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepo := mock.NewMockUserRepo(ctrl)
			server := NewUserServiceServer(mockUserRepo)

			tt.expects(mockUserRepo)
			got, err := server.GetUserList(tt.args.ctx, tt.args.request)
			tt.asserts(t, got, err)
		})
	}
}

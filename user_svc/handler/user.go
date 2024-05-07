package handler

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zycgary/mxshop-go/user_svc/data"
	. "github.com/zycgary/mxshop-go/user_svc/global"
	"github.com/zycgary/mxshop-go/user_svc/model"
	userproto "github.com/zycgary/mxshop-go/user_svc/proto"

	"github.com/anaskhan96/go-password-encoder"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var (
	passwordOptions = &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
)

func UserModelToUserInfo(user *model.User) *userproto.UserInfo {
	info := userproto.UserInfo{
		Id:       uint64(user.ID),
		Nickname: user.Nickname,
		Password: user.Password,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Role:     user.Role,
	}

	if user.Birthday != nil {
		info.Birthday = uint64(user.Birthday.Unix())
	}

	return &info
}

func adjustPaginationParams(page, pageSize int32) (p, ps int32) {
	if page == 0 {
		page = 1
	}
	switch {
	case pageSize == 0:
		pageSize = 10
	case pageSize > 100:
		pageSize = 100
	}

	return page, pageSize
}

var _ userproto.UserServiceServer = (*UserServiceServer)(nil)

type UserServiceServer struct {
	userproto.UnimplementedUserServiceServer
	userRepo data.UserRepo
}

func NewUserServiceServer(ur data.UserRepo) *UserServiceServer {
	return &UserServiceServer{userRepo: ur}
}

func (u *UserServiceServer) GetUserList(_ context.Context, request *userproto.GetUserListRequest) (*userproto.UserListResponse, error) {
	zap.S().Debug("getting user list")
	var response userproto.UserListResponse

	p, ps := adjustPaginationParams(int32(request.Page), int32(request.PageSize))
	total, users, err := u.userRepo.ListUser(p, ps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user list")
	}

	response.Total = total
	response.Data = make([]*userproto.UserInfo, 0, ps)

	for _, user := range users {
		userInfo := UserModelToUserInfo(user)
		response.Data = append(response.Data, userInfo)
	}

	return &response, nil
}

func (u *UserServiceServer) GetUserById(_ context.Context, request *userproto.IdRequest) (*userproto.UserInfoResponse, error) {
	response := &userproto.UserInfoResponse{}

	var user model.User
	r := DB.First(&user, request.Id)
	if err := r.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		return nil, err
	}

	userInfo := UserModelToUserInfo(&user)
	response.Data = userInfo

	return response, nil
}

func (u *UserServiceServer) GetUserByMobile(_ context.Context, request *userproto.MobileRequest) (*userproto.UserInfoResponse, error) {
	response := &userproto.UserInfoResponse{}

	var user model.User
	r := DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if err := r.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		return nil, err
	}

	userInfo := UserModelToUserInfo(&user)
	response.Data = userInfo

	return response, nil
}

func (u *UserServiceServer) CreateUser(_ context.Context, request *userproto.CreateUserRequest) (*userproto.UserInfoResponse, error) {
	response := &userproto.UserInfoResponse{}

	// Check if user with the same mobile number already exists
	var existingUserCount int64
	r := DB.Where(&model.User{Mobile: request.Mobile}).Count(&existingUserCount)
	if err := r.Error; err != nil {
		return nil, err
	}
	if existingUserCount > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "Mobile number is already in used.")
	}

	var user = model.User{}
	user.Mobile = request.Mobile
	user.Nickname = request.Nickname

	// Encrypt password
	salt, encoded := password.Encode(request.Password, passwordOptions)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encoded)

	// Create user
	r = DB.Create(&user)
	if r.Error != nil {
		return nil, status.Errorf(codes.Internal, r.Error.Error())
	}

	response.Data = UserModelToUserInfo(&user)

	return response, nil
}

func (u *UserServiceServer) UpdateUser(_ context.Context, request *userproto.UpdateUserRequest) (*emptypb.Empty, error) {
	var user model.User
	if r := DB.First(&user, request.Id); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		return nil, r.Error
	}

	birthday := time.Unix(int64(request.Birthday), 0)
	user.Nickname = request.Nickname
	user.Birthday = &birthday
	user.Gender = request.Gender

	if r := DB.Save(user); r.Error != nil {
		return nil, status.Errorf(codes.Internal, r.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u *UserServiceServer) CheckPassword(_ context.Context, request *userproto.CheckPasswordRequest) (*userproto.CheckPasswordResponse, error) {
	response := &userproto.CheckPasswordResponse{}
	passwordInfo := strings.Split(request.EncryptedPassword, "$")
	check := password.Verify(request.Password, passwordInfo[2], passwordInfo[3], passwordOptions)
	response.Success = check

	return response, nil
}

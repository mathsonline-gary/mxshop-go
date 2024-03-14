package user_controller

import (
	"fmt"
	"net/http"
	"time"

	"mxshop-go/user_api/global/response"
	"mxshop-go/user_api/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ip   = "localhost"
	port = 50051
)

func GrpcErrorToHttpResponse(err error, ctx *gin.Context) {
	stat := http.StatusInternalServerError
	msg := "Internal error"

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				stat = http.StatusNotFound
				msg = "User list not found"
			case codes.InvalidArgument:
				stat = http.StatusUnprocessableEntity
				msg = "Invalid request"
			case codes.Internal:
			default:
				stat = http.StatusInternalServerError
				msg = "Internal error"
			}
		}

		ctx.JSON(stat, gin.H{
			"message": msg,
		})
	}
}

func Index(ctx *gin.Context) {
	ucc, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[User][Index] failed to connect to service")
		GrpcErrorToHttpResponse(err, ctx)
		return
	}

	userServiceClient := proto.NewUserServiceClient(ucc)
	rsp, err := userServiceClient.GetUserList(ctx, &proto.GetUserListRequest{
		Page:     0,
		PageSize: 0,
	})
	if err != nil {
		zap.S().Errorw(
			"[User][Index] failed to get user list",
			"error", err,
		)
		GrpcErrorToHttpResponse(err, ctx)
		return
	}

	data := make([]response.UserResponse, 0)
	for _, v := range rsp.Data {
		data = append(data, response.UserResponse{
			ID:       v.Id,
			Nickname: v.Nickname,
			Birthday: response.Birthday(time.Unix(int64(v.Birthday), 0)),
			Gender:   v.Gender,
			Mobile:   v.Mobile,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
		"meta": map[string]interface{}{
			"total": rsp.Total,
		},
	})
	return
}

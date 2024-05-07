package initialize

import (
	"fmt"

	"github.com/zycgary/mxshop-go/user_api/global"
	"github.com/zycgary/mxshop-go/user_api/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserSvcClient() {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.Config.ConsulConfig.Host, global.Config.ConsulConfig.Port, global.Config.UserSvcConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("connect to [user grpc service] failed")
	}

	global.UserSvcClient = proto.NewUserServiceClient(conn)
}

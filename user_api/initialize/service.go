package initialize

import (
	"fmt"

	"mxshop-go/user_api/global"
	"mxshop-go/user_api/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserSvcConn() {
	// Get user service info from Consul
	userSvcHost := ""
	userSvcPort := 0

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	res, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSvcConfig.Name))
	if err != nil {
		panic(err)
	}
	for _, service := range res {
		userSvcHost = service.Address
		userSvcPort = service.Port
		break
	}

	if userSvcHost == "" {
		zap.S().Fatal("user grpc service unavailable")
		return
	}

	ucc, err := grpc.Dial(fmt.Sprintf("%s:%d", userSvcHost, userSvcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[User][Index] failed to connect to user grpc service")
		return
	}

	global.UserSvcClient = proto.NewUserServiceClient(ucc)
}

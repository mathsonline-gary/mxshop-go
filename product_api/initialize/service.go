package initialize

import (
	"fmt"

	"mxshop-go/product_api/global"
	"mxshop-go/product_svc/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProductSvcClient() {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.Config.ConsulConfig.Host, global.Config.ConsulConfig.Port, global.Config.ProductSvcConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("connect to [product grpc service] failed", err)
	}

	global.ProductSvcClient = proto.NewProductServiceClient(conn)
}

package data

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	usv1 "github.com/zycgary/mxshop-go/api/user/service/v1"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/config"
	"github.com/zycgary/mxshop-go/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Data struct {
	UserServiceClient usv1.UserServiceClient

	logger *log.Sugar
}

func NewData(conf *config.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{
		logger: log.NewSugar(logger),
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", conf.UserService.Host, conf.UserService.Port, conf.UserService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		data.logger.Fatalf("did not connect: %v", err)
	}
	cleanup := func() {
		_ = conn.Close()
	}

	data.UserServiceClient = usv1.NewUserServiceClient(conn)
	return data, cleanup, nil
}

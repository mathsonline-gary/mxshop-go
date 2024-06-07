package server

import (
	"context"
	"fmt"
	"net"

	pbv1 "github.com/zycgary/mxshop-go/api/user/service/v1"
	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/config"
	"github.com/zycgary/mxshop-go/pkg/log"
	"github.com/zycgary/mxshop-go/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var _ transport.Server = (*server)(nil)

type server struct {
	*grpc.Server

	logger  *log.Sugar
	network string
	host    string
	port    int32
}

func NewGRPCServer(conf config.Server, controller pbv1.UserServiceServer, logger log.Logger) transport.Server {
	s := &server{
		Server:  grpc.NewServer(),
		logger:  log.NewSugar(logger),
		network: conf.Network,
		host:    conf.Host,
		port:    conf.Port,
	}

	s.logger.Debugf("Registering gRPC server")

	pbv1.RegisterUserServiceServer(s, controller)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	return s
}

func (s *server) Start(ctx context.Context) error {
	lis, err := net.Listen(s.network, fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return err
	}

	return s.Serve(lis)
}

func (s *server) Stop(ctx context.Context) error {
	s.Server.GracefulStop()

	return nil
}

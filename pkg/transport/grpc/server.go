package grpc

import (
	"context"
	"net/url"

	"github.com/zycgary/mxshop-go/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Option func(*Server)

func WithGRPCOptions() {

}

// Server is a gRPC Server wrapper.
type Server struct {
	*grpc.Server
	grpcOpts []grpc.ServerOption
	health   *health.Server
}

func NewServer() *Server {
	// TODO implement me
	panic("implement me")
}

func (s *Server) Endpoint() (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

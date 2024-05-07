package rpcserver

import (
	"github.com/zycgary/mxshop-go/pkg/host"
	"net"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServerOption func(s *Server)

type Server struct {
	*grpc.Server

	address            string
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	grpcOptions        []grpc.ServerOption
	listener           net.Listener

	healthServer *health.Server
	endpoint     *url.URL
}

func (s *Server) setEndpoint() error {
	if s.listener == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.listener = lis
	}

	addr, err := host.Extract(s.address, s.listener)
	if err != nil {
		_ = s.listener.Close()

		return err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}

	return nil
}

func NewServer(options ...ServerOption) *Server {
	srv := &Server{
		address:      ":0",
		healthServer: health.NewServer(),
	}

	for _, opt := range options {
		opt(srv)
	}

	srv.setEndpoint()

	// TODO: add default interceptors
	// ...

	grpcOptions := grpc.ChainUnaryInterceptor(srv.unaryInterceptors...)
	if srv.grpcOptions != nil && len(srv.grpcOptions) > 0 {
		srv.grpcOptions = append(srv.grpcOptions, grpcOptions)
	}

	srv.Server = grpc.NewServer(srv.grpcOptions...)

	// Extract endpoint from address
	err := srv.setEndpoint()
	if err != nil {
		return nil
	}

	// Register health server
	if srv.healthServer != nil {
		grpc_health_v1.RegisterHealthServer(srv.Server, srv.healthServer)
	}

	return srv
}

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func WithUnaryInterceptors(interceptors []grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInterceptors = interceptors
	}
}

func WithStreamInterceptors(interceptors []grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInterceptors = interceptors
	}
}

func WithGRPCOptions(opts []grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOptions = opts
	}
}

func WithListener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.listener = lis
	}
}

func WithHealthServer(srv *health.Server) ServerOption {
	return func(s *Server) {
		s.healthServer = srv
	}
}

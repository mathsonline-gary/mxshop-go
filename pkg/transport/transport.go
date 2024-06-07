package transport

import (
	"context"
	"net/url"

	"google.golang.org/grpc"
)

// Server is transport server.
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Endpointer is registry endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}

// GRPCServerGetter returns the (cloned) v1 server instance.
type GRPCServerGetter interface {
	GetGRPCServer() *grpc.Server
}

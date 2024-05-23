package app

import (
	"context"
	"net/url"
	"os"

	"github.com/zycgary/mxshop-go/pkg/log"
	"google.golang.org/grpc"
)

type Option func(o *options)

type options struct {
	id       string
	name     string
	endpoint *url.URL

	ctx     context.Context
	signals []os.Signal

	logger     log.Logger
	grpcServer *grpc.Server

	// Before and After funcs
	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// WithID sets the app service id.
func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

// WithName sets the app service name.
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithEndpoint sets the app service endpoint.
func WithEndpoint(endpoint *url.URL) Option {
	return func(o *options) {
		o.endpoint = endpoint
	}
}

// WithContext sets the context, and will be used to derive context of App struct in New().
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithSignals sets the app signals.
func WithSignals(signals ...os.Signal) Option {
	return func(o *options) {
		o.signals = signals
	}
}

// WithLogger sets the app logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithGRPCServer sets the app grpc server.
func WithGRPCServer(s *grpc.Server) Option {
	return func(o *options) {
		o.grpcServer = s
	}
}

// WithBeforeStart BeforeStart run funcs before app starts
func WithBeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

// WithBeforeStop BeforeStop run funcs before app stops
func WithBeforeStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

// WithAfterStart AfterStart run funcs after app starts
func WithAfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

// WithAfterStop AfterStop run funcs after app stops
func WithAfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}

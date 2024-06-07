package app

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/zycgary/mxshop-go/pkg/log"
	"github.com/zycgary/mxshop-go/pkg/registry"
	"github.com/zycgary/mxshop-go/pkg/transport"
)

type Option func(o *options)

type options struct {
	id       string
	name     string
	tags     []string
	metadata map[string]string
	endpoint *url.URL

	ctx     context.Context
	signals []os.Signal

	logger log.Logger

	servers     []transport.Server
	stopTimeout time.Duration

	registrar        registry.Registrar
	registrarTimeout time.Duration

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

// WithTags sets the app service tags.
func WithTags(tags ...string) Option {
	return func(o *options) {
		o.tags = tags
	}
}

// WithMetadata sets the app service metadata.
func WithMetadata(md map[string]string) Option {
	return func(o *options) {
		o.metadata = md
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

// WithServers sets the app servers.
func WithServers(srv ...transport.Server) Option {
	return func(o *options) {
		o.servers = srv
	}
}

// WithRegistrar sets the app registrar.
func WithRegistrar(r registry.Registrar) Option {
	return func(o *options) {
		o.registrar = r
	}
}

// WithRegistrarTimeout sets the app registrar timeout.
func WithRegistrarTimeout(t time.Duration) Option {
	return func(o *options) {
		o.registrarTimeout = t
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

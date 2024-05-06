package app

import (
	"net/url"
	"os"
	"syscall"
	"time"

	"mxshop-go/gmicro/mesh"

	"github.com/google/uuid"
)

type Option func(*options)

type options struct {
	id        string
	name      string
	endpoints []url.URL

	exitSignals []os.Signal

	registry          mesh.Registry
	registerTimeout   time.Duration
	deregisterTimeout time.Duration
}

func NewDefaultOptions() *options {
	opts := options{
		exitSignals: []os.Signal{
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
		},

		registerTimeout:   10 * time.Second,
		deregisterTimeout: 10 * time.Second,
	}

	if id, err := uuid.NewUUID(); err == nil {
		opts.id = id.String()
	}

	return &opts
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}
func WithEndpoints(endpoints []url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

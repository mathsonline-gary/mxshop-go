package app

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/zycgary/mxshop-go/pkg/registry"
	"golang.org/x/sync/errgroup"
)

type App struct {
	opts     options
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	instance *registry.Instance
}

func New(opts ...Option) *App {
	o := options{
		ctx:              context.Background(),
		signals:          []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
	}

	// generate default id
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancel := context.WithCancel(o.ctx)

	return &App{
		opts:   o,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run starts the app and waits for the stop signal.
// It includes:
// - start server(s)
// - stop server(s) when error occurs
// - stop the app gracefully when stop signal received
func (a *App) Run() error {
	// Build instance.
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()

	for _, fn := range a.opts.beforeStart {
		if err := fn(a.ctx); err != nil {
			return err
		}
	}

	lis, err := net.Listen("tcp", a.opts.endpoint.String())
	if err != nil {
		return err
	}
	eg, ectx := errgroup.WithContext(a.ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)

	// stop server when App context done
	eg.Go(func() error {
		<-ectx.Done() // wait for errgroup cancel signal
		a.opts.grpcServer.Stop()
		return nil
	})

	// Start server(s)
	eg.Go(func() error {
		wg.Done() // defer is not needed here, as it is to identify the server has begun to start before the registration, not to identify the server has started.
		log.Printf("server listening at %v", lis.Addr())
		if err := a.opts.grpcServer.Serve(lis); err != nil {
			return err
		}
		return nil
	})
	wg.Wait() // wait for server to start

	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(a.ctx, a.opts.registrarTimeout)
		defer rcancel()
		if err = a.opts.registrar.Register(rctx, instance); err != nil {
			return err
		}
	}

	for _, fn := range a.opts.afterStart {
		if err = fn(a.ctx); err != nil {
			return err
		}
	}

	// Wait for stop signals, and stop the app gracefully.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, a.opts.signals...)
	eg.Go(func() error {
		select {
		case <-ectx.Done():
			return nil
		case <-sc:
			return a.Stop()
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	for _, fn := range a.opts.afterStop {
		if err := fn(a.ctx); err != nil {
			return err
		}
	}

	return nil
}

// Stop stops the app gracefully.
func (a *App) Stop() error {
	// Deregister service from consul.
	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()
	if a.opts.registrar != nil && instance != nil {
		ctx, cancel := context.WithTimeout(a.ctx, a.opts.registrarTimeout)
		defer cancel()
		if err := a.opts.registrar.Deregister(ctx, instance); err != nil {
			return err
		}
	}

	// Close logger.
	if a.opts.logger != nil {
		_ = a.opts.logger.Close()
	}

	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func (a *App) buildInstance() (*registry.Instance, error) {
	return &registry.Instance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Tags:      a.opts.tags,
		Metadata:  a.opts.metadata,
		Endpoints: []string{a.opts.endpoint.String()},
	}, nil
}

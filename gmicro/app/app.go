package app

import (
	"context"
	"mxshop-go/gmicro/service_manager"
	"os"
	"os/signal"
	"sync"
)

type App struct {
	options  options
	mu       sync.Mutex
	instance *service_manager.ServiceInstance
}

func New(opts ...Option) *App {
	options := NewDefaultOptions()

	for _, opt := range opts {
		opt(options)
	}
	return &App{
		options: *options,
	}
}

func (app *App) getInstance() *service_manager.ServiceInstance {
	app.mu.Lock()
	ins := app.instance
	app.mu.Unlock()

	return ins
}

func (app *App) setInstance(instance *service_manager.ServiceInstance) {
	app.mu.Lock()
	app.instance = instance
	app.mu.Unlock()
}

func (app *App) Run() error {
	// Create service instance
	si, err := app.buildServiceInstance()
	if err != nil {
		return err
	}
	app.setInstance(si)

	// Register service
	if app.options.registry != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), app.options.registerTimeout)
		defer rcancel()
		if err := app.options.registry.Register(rctx, si); err != nil {
			return err
		}
	}

	// Monitor exit signals, this section should be placed at the bottom of the method
	ech := make(chan os.Signal, 1)
	signal.Notify(ech, app.options.exitSignals...)
	<-ech

	return nil
}

func (app *App) Stop() error {
	ins := app.getInstance()

	if app.options.registry != nil && ins != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), app.options.deregisterTimeout)
		defer rcancel()
		if err := app.options.registry.Deregister(rctx, ins); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) buildServiceInstance() (*service_manager.ServiceInstance, error) {
	return &service_manager.ServiceInstance{
		ID:        app.options.id,
		Name:      app.options.name,
		Endpoints: app.options.endpoints,
	}, nil
}

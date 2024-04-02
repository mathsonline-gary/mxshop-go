package app

import (
	"mxshop-go/gmicro/registery"
)

type App struct {
	options options
}

func New(opts ...Option) *App {
	options := options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &App{
		options: options,
	}
}

func (app *App) Run() error {
	return nil
}

func (app *App) Stop() error {
	return nil
}

func (app *App) buildServiceInstance() (*registery.ServiceInstance, error) {
	return &registery.ServiceInstance{
		ID:        app.options.id,
		Name:      app.options.name,
		Endpoints: app.options.endpoints,
	}, nil
}

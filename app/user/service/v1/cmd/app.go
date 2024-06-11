package main

import (
	"fmt"
	"net/url"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/config"
	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/data"
	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/logic"
	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/server"
	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/service"
	"github.com/zycgary/mxshop-go/pkg/app"
	zaplog "github.com/zycgary/mxshop-go/pkg/log/zap"
	"github.com/zycgary/mxshop-go/pkg/registry/consul"
	"go.uber.org/zap"
)

func newApp(conf config.Config) (*app.App, error) {
	// Build endpoint.
	endpoint, err := url.Parse(fmt.Sprintf("grpc://%s:%d", conf.App.Host, conf.App.Port))
	if err != nil {
		return nil, err
	}

	// Initialize logger.
	var l *zap.Logger
	if conf.App.Env == "production" {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	logger := zaplog.NewLogger(l)

	// Initialize DB.
	db, _ := data.NewDB(conf.DB, logger)

	// Initialize GRPC server.
	repo := data.NewUserRepository(db, logger)
	uc := logic.NewUserUseCase(repo, logger)
	svc := service.NewUserService(uc, logger)
	s := server.NewGRPCServer(conf.Server, svc, logger)

	// Initialize service registrar.
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, err
	}
	check := &consulapi.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", conf.Consul.Service.Check.Endpoint, conf.App.Port),
		Timeout:                        fmt.Sprintf("%ds", conf.Consul.Service.Check.Timeout),
		Interval:                       fmt.Sprintf("%ds", conf.Consul.Service.Check.Interval),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%dm", conf.Consul.Service.Check.DeregisterAfter),
	}
	registrar := consul.New(
		client,
		consul.WithCheck(check),
	)

	return app.New(
		app.WithName(conf.App.Name),
		app.WithTags(conf.Consul.Service.Tags...),
		app.WithMetadata(map[string]string{}),
		app.WithEndpoint(endpoint),
		app.WithLogger(logger),
		app.WithServers(s),
		app.WithRegistrar(registrar),
	), nil
}

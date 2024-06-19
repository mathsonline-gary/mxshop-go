package main

import (
	"fmt"
	"net/url"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/config"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/data"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/logic"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/server/http"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/service"
	"github.com/zycgary/mxshop-go/pkg/app"
	zaplog "github.com/zycgary/mxshop-go/pkg/log/zap"
	"github.com/zycgary/mxshop-go/pkg/registry/consul"
	"go.uber.org/zap"
)

func newApp(conf *config.Config) (*app.App, func(), error) {
	// Build endpoint.
	endpoint, err := url.Parse(fmt.Sprintf("%s://%s:%d", conf.Server.HTTP.Scheme, conf.Server.HTTP.Host, conf.Server.HTTP.Port))
	if err != nil {
		return nil, nil, err
	}

	// Initialize logger.
	var l *zap.Logger
	if conf.App.Env == "production" {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	logger := zaplog.NewLogger(l)

	// Initialize data.
	d, cleanup, err := data.NewData(conf.Data, logger)
	if err != nil {
		return nil, cleanup, err
	}

	// Initialize GRPC server.
	repo := data.NewUserRepository(d.UserServiceClient, logger)
	uuc := logic.NewUserUseCase(repo, logger)
	auc := logic.NewAuthUseCase(conf.Auth.Secret, repo, logger)
	us := service.NewUserService(uuc, logger)
	as := service.NewAuthService(auc, logger)
	s := http.NewHttpServer(conf, us, as, logger)

	// Initialize service registrar.
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, cleanup, err
	}
	check := &consulapi.AgentServiceCheck{
		HTTP:                           conf.Registry.Check.Endpoint,
		Timeout:                        fmt.Sprintf("%ds", conf.Registry.Check.Timeout),
		Interval:                       fmt.Sprintf("%ds", conf.Registry.Check.Interval),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%dm", conf.Registry.Check.DeregisterAfter/60),
	}
	registrar := consul.New(
		client,
		consul.WithCheck(check),
	)

	return app.New(
		app.WithName(conf.Registry.Name),
		app.WithTags(conf.Registry.Tags...),
		app.WithMetadata(map[string]string{}),
		app.WithEndpoint(endpoint),
		app.WithLogger(logger),
		app.WithServers(s),
		app.WithRegistrar(registrar),
	), cleanup, nil
}

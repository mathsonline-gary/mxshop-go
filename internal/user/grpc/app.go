package grpc

import (
	"fmt"
	"net/url"

	consulapi "github.com/hashicorp/consul/api"
	upbv1 "github.com/zycgary/mxshop-go/api/user/v1"
	"github.com/zycgary/mxshop-go/internal/user/grpc/config"
	ucv1 "github.com/zycgary/mxshop-go/internal/user/grpc/controller/v1"
	urv1 "github.com/zycgary/mxshop-go/internal/user/grpc/repository/v1"
	usv1 "github.com/zycgary/mxshop-go/internal/user/grpc/service/v1"
	"github.com/zycgary/mxshop-go/pkg/app"
	zaplog "github.com/zycgary/mxshop-go/pkg/log/zap"
	"github.com/zycgary/mxshop-go/pkg/registry/consul"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewApp(conf config.Config) (*app.App, error) {
	// Build endpoint.
	endpoint, err := url.Parse(fmt.Sprintf("%s:%d", conf.App.Host, conf.App.Port))
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
	db, _ := urv1.NewDB(conf.DB, logger)

	// Initialize GRPC server.
	repo := urv1.NewUserRepository(db, logger)
	service := usv1.NewUserService(repo, logger)
	controller := ucv1.NewUserController(service, logger)
	s := grpc.NewServer()
	upbv1.RegisterUserServiceServer(s, controller)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

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
		app.WithGRPCServer(s),
		app.WithRegistrar(registrar),
	), nil
}

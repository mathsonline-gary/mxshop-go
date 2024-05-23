package main

import (
	"flag"
	"fmt"
	"net/url"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/zycgary/mxshop-go/order_svc/config"
	"github.com/zycgary/mxshop-go/order_svc/data"
	"github.com/zycgary/mxshop-go/order_svc/handler"
	orderproto "github.com/zycgary/mxshop-go/order_svc/proto"
	"github.com/zycgary/mxshop-go/pkg/app"
	zaplog "github.com/zycgary/mxshop-go/pkg/log/zap"
	"github.com/zycgary/mxshop-go/pkg/registry/consul"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	env = flag.String("env", "local", "The running environment of the service")
)

func newApp(conf config.Config) (*app.App, error) {
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

	defer func(logger *zaplog.Logger) {
		_ = logger.Close()
	}(logger)

	// Initialize DB.
	db, _ := data.NewGormDB(conf.DB, logger)

	// Initialize order service.
	orderRepo := data.NewOrderRepo(db)
	orderService := handler.NewOrderService(
		handler.WithRepo(orderRepo),
	)

	// Initialize GRPC server.
	s := grpc.NewServer()
	orderproto.RegisterOrderServiceServer(s, orderService)
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

func main() {
	flag.Parse()

	// Load config.
	var conf config.Config
	if err := conf.Load("config/order", fmt.Sprintf("grpc.%s", *env), "yaml"); err != nil {
		panic(err)
	}

	a, err := newApp(conf)
	if err != nil {
		panic(err)
	}

	// Start and wait for stop signal.
	if err := a.Run(); err != nil {
		panic(err)
	}
}

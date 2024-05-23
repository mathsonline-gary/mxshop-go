package main

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/zycgary/mxshop-go/order_svc/config"
	"github.com/zycgary/mxshop-go/order_svc/data"
	"github.com/zycgary/mxshop-go/order_svc/handler"
	orderproto "github.com/zycgary/mxshop-go/order_svc/proto"
	"github.com/zycgary/mxshop-go/pkg/app"
	zaplog "github.com/zycgary/mxshop-go/pkg/log/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	env = flag.String("env", "local", "The running environment of the service")
)

func newApp(conf config.Config) (*app.App, error) {
	// Build endpoint
	endpoint, err := url.Parse(fmt.Sprintf("%s:%d", conf.App.Host, conf.App.Port))
	if err != nil {
		return nil, err
	}

	// Initialize logger
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

	// Initialize DB
	db, _ := data.NewGormDB(conf.DB, logger)

	// Initialize order service
	orderRepo := data.NewOrderRepo(db)
	orderService := handler.NewOrderService(
		handler.WithRepo(orderRepo),
	)

	s := grpc.NewServer()
	orderproto.RegisterOrderServiceServer(s, orderService)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	return app.New(
		app.WithName(conf.App.Name),
		app.WithEndpoint(endpoint),
		app.WithLogger(logger),
		app.WithGRPCServer(s),
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

	// Run app
	// Register service to consul
	//client, serviceID, err := registerConsulService(conf)
	//if err != nil {
	//	panic(err)
	//}
	//
	//lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.App.Host, conf.App.Port))
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//
	//// Start service
	//go func() {
	//	log.Printf("server listening at %v", lis.Addr())
	//	if err := s.Serve(lis); err != nil {
	//		log.Fatalf("failed to serve: %v", err)
	//	}
	//}()
	//
	//// Receive quit signal
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//if err := client.Agent().ServiceDeregister(serviceID); err != nil {
	//	zap.S().Info("deregister order service failed")
	//}
	//zap.S().Info("deregister order service successfully")
}

//func registerConsulService(conf config.Config) (client *consulAPI.Client, serviceID string, err error) {
//	c := consulAPI.DefaultConfig()
//	c.Address = fmt.Sprintf("%s:%d", conf.Consul.Host, conf.Consul.Port)
//	client, err = consulAPI.NewClient(c)
//	if err != nil {
//		return nil, "", err
//	}
//	serviceID, _ = uuid.GenerateUUID()
//	registration := &consulAPI.AgentServiceRegistration{
//		Name:    conf.App.Name,
//		ID:      serviceID,
//		Tags:    global.Config.Consul.Service.Tags,
//		Address: conf.App.Host,
//		Port:    int(conf.App.Port),
//		Check: &consulAPI.AgentServiceCheck{
//			GRPC:                           fmt.Sprintf("%s:%d", conf.Consul.Service.Check.Endpoint, conf.App.Port),
//			Timeout:                        fmt.Sprintf("%ds", conf.Consul.Service.Check.Timeout),
//			Interval:                       fmt.Sprintf("%ds", conf.Consul.Service.Check.Interval),
//			DeregisterCriticalServiceAfter: fmt.Sprintf("%dm", conf.Consul.Service.Check.DeregisterAfter),
//		},
//	}
//
//	if err := client.Agent().ServiceRegister(registration); err != nil {
//		return nil, "", err
//	}
//
//	return client, serviceID, nil
//}

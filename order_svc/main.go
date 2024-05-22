package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	consulAPI "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"github.com/zycgary/mxshop-go/order_svc/config"
	"github.com/zycgary/mxshop-go/order_svc/data"
	"github.com/zycgary/mxshop-go/order_svc/global"
	"github.com/zycgary/mxshop-go/order_svc/handler"
	orderproto "github.com/zycgary/mxshop-go/order_svc/proto"
	zaplog "github.com/zycgary/mxshop-go/pkg/log/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	env = flag.String("env", "local", "The running environment of the service")
)

func main() {
	flag.Parse()

	// Initialize config
	var conf config.Config
	if err := conf.Load("config/order", fmt.Sprintf("grpc.%s", *env), "yaml"); err != nil {
		panic(err)
	}
	conf.Watch()

	// Initialize logger
	var l *zap.Logger
	if conf.App.Env == "production" {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	logger := zaplog.NewLogger(l)
	defer func(logger *zaplog.Logger) {
		_ = logger.Sync()
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
	log.Printf("%s:%d", conf.App.Host, conf.App.Port)
	client, serviceID, err := registerConsulService(conf)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.App.Host, conf.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start service
	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Receive quit signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("deregister order service failed")
	}
	zap.S().Info("deregister order service successfully")
}

func registerConsulService(conf config.Config) (client *consulAPI.Client, serviceID string, err error) {
	config := consulAPI.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", conf.Consul.Host, conf.Consul.Port)
	client, err = consulAPI.NewClient(config)
	if err != nil {
		return nil, "", err
	}
	serviceID, _ = uuid.GenerateUUID()
	registration := &consulAPI.AgentServiceRegistration{
		Name:    conf.App.Name,
		ID:      serviceID,
		Tags:    global.Config.Consul.Service.Tags,
		Address: conf.App.Host,
		Port:    int(conf.App.Port),
		Check: &consulAPI.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", conf.Consul.Service.Check.Endpoint, conf.App.Port),
			Timeout:                        fmt.Sprintf("%ds", conf.Consul.Service.Check.Timeout),
			Interval:                       fmt.Sprintf("%ds", conf.Consul.Service.Check.Interval),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%dm", conf.Consul.Service.Check.DeregisterAfter),
		},
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return nil, "", err
	}

	return client, serviceID, nil
}

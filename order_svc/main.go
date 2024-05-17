package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zycgary/mxshop-go/order_svc/config"
	"github.com/zycgary/mxshop-go/order_svc/data"
	"github.com/zycgary/mxshop-go/order_svc/global"
	"github.com/zycgary/mxshop-go/order_svc/handler"
	"github.com/zycgary/mxshop-go/order_svc/initialize"
	orderproto "github.com/zycgary/mxshop-go/order_svc/proto"
	"gorm.io/gorm/logger"

	consulAPI "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	initialize.Init()
	var (
		env  = flag.String("env", "local", "The running environment of the service")
		ip   = flag.String("ip", global.Config.App.Host, "The user service IP")
		port = flag.Int("port", int(global.Config.App.Port), "The user service port")
	)
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize config
	var conf config.Config
	if err := conf.Load("config/order", fmt.Sprintf("grpc.%s", *env), "yaml"); err != nil {
		panic(err)
	}
	conf.Watch()

	// Initialize DB
	logLevel := logger.Silent
	if conf.App.Debug {
		logLevel = logger.Info
	}
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLevel,    // Log level
		},
	)
	db := data.NewGormDB(conf.DB, dbLogger)

	// Initialize order service
	orderRepo := data.NewOrderRepo(db)
	orderServiceServer := handler.NewOrderServiceServer(
		handler.WithRepo(orderRepo),
	)

	s := grpc.NewServer()
	orderproto.RegisterOrderServiceServer(s, orderServiceServer)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	client, serviceID, err := registerConsulService(*ip, *port)
	if err != nil {
		panic(err)
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

func registerConsulService(addr string, port int) (client *consulAPI.Client, serviceID string, err error) {
	config := consulAPI.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.Config.Consul.Host, global.Config.Consul.Port)
	client, err = consulAPI.NewClient(config)
	if err != nil {
		return nil, "", err
	}
	serviceID, _ = uuid.GenerateUUID()
	registration := &consulAPI.AgentServiceRegistration{
		Name:    global.Config.App.Name,
		ID:      serviceID,
		Tags:    global.Config.Consul.Service.Tags,
		Address: addr,
		Port:    port,
		Check: &consulAPI.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", global.Config.Consul.Service.Check.Endpoint, port),
			Timeout:                        fmt.Sprintf("%ds", global.Config.Consul.Service.Check.Timeout),
			Interval:                       fmt.Sprintf("%ds", global.Config.Consul.Service.Check.Interval),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%dm", global.Config.Consul.Service.Check.DeregisterAfter),
		},
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return nil, "", err
	}

	return client, serviceID, nil
}

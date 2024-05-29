package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/zycgary/mxshop-go/stock_svc/global"
	"github.com/zycgary/mxshop-go/stock_svc/handler"
	"github.com/zycgary/mxshop-go/stock_svc/initialize"
	stockproto "github.com/zycgary/mxshop-go/stock_svc/proto"

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
		ip   = flag.String("ip", global.Config.AppConfig.Host, "The user service IP")
		port = flag.Int("port", global.Config.AppConfig.Port, "The user service port")
	)
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	stockproto.RegisterStockServiceServer(s, &handler.StockServiceServer{})
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
		zap.S().Info("deregister stock service failed")
	}
	zap.S().Info("deregister stock service successfully")
}

func registerConsulService(addr string, port int) (client *consulAPI.Client, serviceID string, err error) {
	config := consulAPI.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.Config.ConsulConfig.Host, global.Config.ConsulConfig.Port)
	client, err = consulAPI.NewClient(config)
	if err != nil {
		return nil, "", err
	}
	serviceID, _ = uuid.GenerateUUID()
	registration := &consulAPI.AgentServiceRegistration{
		Name:    global.Config.AppConfig.Name,
		ID:      serviceID,
		Tags:    []string{"mxshop", "stock", "grpc"},
		Address: addr,
		Port:    port,
		Check: &consulAPI.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", global.Config.ConsulConfig.StockSvc.Check.Host, port),
			Timeout:                        "5s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}
	if err := client.Agent().ServiceRegister(registration); err != nil {
		return nil, "", err
	}

	return client, serviceID, nil
}

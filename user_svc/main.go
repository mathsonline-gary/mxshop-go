package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop-go/user_svc/data"
	"mxshop-go/user_svc/global"
	"mxshop-go/user_svc/handler"
	"mxshop-go/user_svc/initialize"
	userproto "mxshop-go/user_svc/proto"

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
	userRepo := data.NewUserRepo(global.DB)
	userproto.RegisterUserServiceServer(s, handler.NewUserServiceServer(userRepo))
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
		zap.S().Info("deregister user service failed")
	}
	zap.S().Info("deregister user service successfully")
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
		Tags:    []string{"mxshop", "user", "svc"},
		Address: addr,
		Port:    port,
		Check: &consulAPI.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", global.Config.ConsulConfig.UserSvc.Check.Host, port),
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

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"mxshop-go/user_svc/global"
	"mxshop-go/user_svc/handler"
	"mxshop-go/user_svc/initialize"
	userproto "mxshop-go/user_svc/proto"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	initialize.Init()

	var (
		ip   = flag.String("ip", global.ServerConfig.AppConfig.Host, "The user service IP")
		port = flag.Int("port", global.ServerConfig.AppConfig.Port, "The user service port")
	)
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userproto.RegisterUserServiceServer(s, &handler.UserServiceServer{})
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	//Register service in Consul
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	serviceID, _ := uuid.GenerateUUID()
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.AppConfig.Name,
		ID: serviceID,
		Tags:    []string{"mxshop", "user", "svc"},
		Address: *ip,
		Port:    *port,
		Check: &api.AgentServiceCheck{
			GRPC:     fmt.Sprintf("%s:%d", *ip, *port),
			Timeout:  "5s",
			Interval: "10s",
			//DeregisterCriticalServiceAfter: "30s",
		},
	}
	if err := client.Agent().ServiceRegister(registration); err != nil {
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
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("deregister user service failed")
	}
	zap.S().Info("deregister user service successfully")
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"mxshop-go/user_svc/global"
	"mxshop-go/user_svc/handler"
	"mxshop-go/user_svc/initialize"
	userproto "mxshop-go/user_svc/proto"

	"google.golang.org/grpc"
)

func main() {
	initialize.Init()

	var (
		ip   = flag.String("ip", global.ServerConfig.AppConfig.Host, "The user service IP")
		port = flag.Int("port", global.ServerConfig.AppConfig.Post, "The user service port")
	)
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userproto.RegisterUserServiceServer(s, &handler.UserServiceServer{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"mxshop-go/user_svc/handler"
	userproto "mxshop-go/user_svc/proto"

	"google.golang.org/grpc"
)

var (
	ip = flag.String("ip", "localhost", "The user service IP")
	port = flag.Int("port", 50051, "The user service port")
)

func main() {
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

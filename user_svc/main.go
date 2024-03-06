package main

import (
	"fmt"
	"net"
	"net/rpc"

	"mxshop-go/user_svc/handler"
	"mxshop-go/user_svc/server_stub"
)

func main() {
	_ = server_stub.RegisterHelloService(&handler.HelloService{})

	lis, _ := net.Listen("tcp", ":9090")
	fmt.Println("user_svc is listening...")

	for {
		conn, _ := lis.Accept()
		go rpc.ServeConn(conn)
	}
}

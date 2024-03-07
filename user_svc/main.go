package main

import (
	"fmt"
	"net"
	"net/rpc"
)

func main() {
	lis, _ := net.Listen("tcp", ":9090")
	fmt.Println("user_svc is listening...")

	for {
		conn, _ := lis.Accept()
		go rpc.ServeConn(conn)
	}
}

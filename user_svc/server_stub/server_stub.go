package server_stub

import "net/rpc"

type IHelloService interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(svc IHelloService) error {
	return rpc.RegisterName("HelloService", svc)
}

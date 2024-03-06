package client_stub

import (
	"net/rpc"
)

type HelloServiceStub struct {
	*rpc.Client
}

func NewHelloServiceClient(protocol, addr string) HelloServiceStub {
	conn, _ := rpc.Dial(protocol, addr)
	return HelloServiceStub{conn}
}

func (c *HelloServiceStub) Hello(request string, reply *string) error {
	return c.Call("HelloService.Hello", request, reply)
}

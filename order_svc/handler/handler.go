package handler

import "github.com/zycgary/mxshop-go/order_svc/proto"

var _ proto.OrderServiceServer = (*OrderServiceServer)(nil)

type OrderServiceServer struct {
	proto.UnimplementedOrderServiceServer
}

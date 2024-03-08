package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"mxshop-go/user_svc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn       *grpc.ClientConn
	userClient proto.UserServiceClient
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	userClient = proto.NewUserServiceClient(conn)
}

func TestGetUserList(t *testing.T) {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	userClient = proto.NewUserServiceClient(conn)
	ctx := context.Background()
	response, err := userClient.GetUserList(ctx, &proto.GetUserListRequest{
		Page:     1,
		PageSize: 0,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println(response)
}

func TestGetUserByMobile(t *testing.T) {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	userClient = proto.NewUserServiceClient(conn)
	ctx := context.Background()
	response, err := userClient.GetUserByMobile(ctx, &proto.MobileRequest{
		Mobile: "0411111111",
	})
	if err != nil {
		fmt.Println(err)
	}
	log.Println(response)
}

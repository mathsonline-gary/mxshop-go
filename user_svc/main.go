package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"log"
	"net"

	"mxshop-go/user_svc/handler"
	userproto "mxshop-go/user_svc/proto"

	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc"
)

func hashMD5() {
	// Using custom options
	options := &password.Options{
		SaltLen:      10,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	salt, encodedPwd := password.Encode("generic password", options)
	check := password.Verify("generic password", salt, encodedPwd, options)
	fmt.Println(check) // true
	fmt.Println(salt)
	fmt.Println(encodedPwd)
}

var (
	port = flag.Int("port", 50051, "The user service port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
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

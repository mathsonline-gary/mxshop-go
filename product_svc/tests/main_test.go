package tests

import (
	"log"
	"os"
	"testing"

	"mxshop-go/product_svc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn          *grpc.ClientConn
	productClient proto.ProductServiceClient
)

// TestMain is the entry point for testing, and it will run before any test.
func TestMain(m *testing.M) {
	setup()
	code := m.Run() // This will execute all tests
	teardown()
	os.Exit(code)
}

func setup() {
	log.Println("start testing...")
	var err error
	conn, err = grpc.Dial("127.0.0.1:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	productClient = proto.NewProductServiceClient(conn)
}

func teardown() {
	log.Println("end testing...")
	if err := conn.Close(); err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}

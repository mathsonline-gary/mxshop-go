package tests

import (
	"log"
	"os"
	"testing"

	"github.com/zycgary/mxshop-go/stock_svc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn        *grpc.ClientConn
	stockClient proto.StockServiceClient
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
	conn, err = grpc.Dial("127.0.0.1:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	stockClient = proto.NewStockServiceClient(conn)
}

func teardown() {
	log.Println("end testing...")
	if err := conn.Close(); err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}

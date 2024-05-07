package tests

import (
	"context"
	"log"
	"testing"

	"github.com/zycgary/mxshop-go/product_svc/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetAllCategories(t *testing.T) {
	rsp, err := productClient.GetAllCategories(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("total: %d", rsp.Total)
	for _, data := range rsp.Data {
		log.Printf("data: %v", data)
	}
	log.Printf("json data: %s", rsp.JsonData)
}

func TestGetSubCategories(t *testing.T) {
	rsp, err := productClient.GetSubCategories(context.Background(), &proto.GetSubCategoriesRequest{Id: 130364})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("sub-categories: %v", rsp.SubCategories)
}

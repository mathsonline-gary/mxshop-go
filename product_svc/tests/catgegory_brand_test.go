package tests

import (
	"context"
	"log"
	"testing"

	"github.com/zycgary/mxshop-go/product_svc/proto"
)

func TestGetCategoryBrandList(t *testing.T) {
	rsp, err := productClient.GetCategoryBrandList(context.Background(), &proto.GetCategoryBrandListRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, data := range rsp.Data {
		log.Printf("data: %v", data)
	}
}

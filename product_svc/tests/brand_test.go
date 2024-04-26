package tests

import (
	"context"
	"log"
	"testing"

	"mxshop-go/product_svc/proto"
)

func TestGetBrands(t *testing.T) {
	rsp, err := productClient.GetBrands(context.Background(), &proto.GetBrandsRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("total: %d", rsp.Total)
	for _, data := range rsp.Data {
		log.Printf("%d: %s", data.Id, data.Name)
	}
}

func TestCreateBrand(t *testing.T) {
	rsp, err := productClient.CreateBrand(context.Background(), &proto.CreateBrandRequest{
		Name: "test1",
		Logo: "test",
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("id: %d; name: %s; logo: %s", rsp.Data.Id, rsp.Data.Name, rsp.Data.Logo)
}

func TestDeleteBrand(t *testing.T) {
	_, err := productClient.DeleteBrand(context.Background(), &proto.DeleteBrandRequest{
		Id: 1120,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println("Deleted")
}

func TestUpdateBrand(t *testing.T) {
	_, err := productClient.UpdateBrand(context.Background(), &proto.UpdateBrandRequest{
		Id:   1121,
		Name: "test2",
		Logo: "test2",
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println("Updated")
}

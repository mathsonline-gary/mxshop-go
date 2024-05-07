package tests

import (
	"context"
	"log"
	"testing"

	"github.com/zycgary/mxshop-go/product_svc/proto"
)

func TestFilterProducts(t *testing.T) {
	var priceMin int32 = 90
	var Brand int32 = 614
	var Category int32 = 130364
	test := []*proto.FilterProductsRequest{
		{
			PriceMin: &priceMin,
		},
		{
			Brand:       &Brand,
			TopCategory: &Category,
		},
	}

	for _, req := range test {
		rsp, err := productClient.FilterProducts(context.Background(), req)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Printf("total: %d", rsp.Total)
		for _, productInfo := range rsp.Data {
			log.Printf("id: %d; category: %v", productInfo.Id, productInfo.Category)
		}
	}
}

func TestBatchGetProducts(t *testing.T) {
	req := &proto.BatchGetProductsRequest{
		Ids: []int32{421, 422, 843},
	}
	rsp, err := productClient.BatchGetProducts(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, productInfo := range rsp.Data {
		log.Printf("id: %d; category: %v", productInfo.Id, productInfo.Category)
	}
}

func TestCreateProduct(t *testing.T) {
	req := &proto.CreateProductRequest{
		Name:         "test",
		SerialNumber: "SN123456",
		Stocks:       10,
		MarketPrice:  30,
		ShopPrice:    20,
		Brief:        "Test Brief",
		Description:  "Test Desc",
		FreeShipping: false,
		Images:       []string{"https://www.test.com/1.jpg", "https://www.test.com/2.jpg"},
		DescImages:   []string{"https://www.test.com/3.jpg", "https://www.test.com/4.jpg"},
		FrontImage:   "https://www.test.com/3.jpg",
		IsNew:        true,
		IsHot:        false,
		OnSale:       false,
		CategoryId:   130366,
		BrandId:      614,
	}
	rsp, err := productClient.CreateProduct(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("product id: %d", rsp.Data.Id)
}

func TestGetProduct(t *testing.T) {
	req := &proto.GetProductRequest{
		Id: 421,
	}
	rsp, err := productClient.GetProduct(context.Background(), req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("product: %v", rsp.Data)
}

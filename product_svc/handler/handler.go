package handler

import (
	"context"

	"mxshop-go/product_svc/proto"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type ProductServiceServer struct {
	proto.UnimplementedProductServiceServer
}

var _ proto.ProductServiceServer = ProductServiceServer{}

func (p ProductServiceServer) FilterProducts(ctx context.Context, request *proto.FilterProductsRequest) (*proto.FilterProductsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) BatchGetProducts(ctx context.Context, request *proto.BatchGetProductsRequest) (*proto.BatchGetProductsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) CreateProduct(ctx context.Context, request *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) DeleteProduct(ctx context.Context, request *proto.DeleteProductRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) UpdateProduct(ctx context.Context, request *proto.UpdateProductRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) GetProduct(ctx context.Context, request *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) GetAllCategories(ctx context.Context, empty *emptypb.Empty) (*proto.GetAllCategoriesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) GetSubCategories(ctx context.Context, request *proto.GetSubCategoriesRequest) (*proto.GetSubCategoriesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) CreateCategory(ctx context.Context, request *proto.CreateCategoryRequest) (*proto.CreateCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) UpdateCategory(ctx context.Context, request *proto.UpdateCategoryRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) GetCategoryBrandList(ctx context.Context, request *proto.GetCategoryBrandListRequest) (*proto.GetCategoryBrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) GetBrandsByCategory(ctx context.Context, request *proto.GetBrandsByCategoryRequest) (*proto.GetBrandsByCategoryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) CreateCategoryBrand(ctx context.Context, request *proto.CreateCategoryBrandRequest) (*proto.CreateCategoryBrandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) DeleteCategoryBrand(ctx context.Context, request *proto.DeleteCategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProductServiceServer) UpdateCategoryBrand(ctx context.Context, request *proto.UpdateCategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

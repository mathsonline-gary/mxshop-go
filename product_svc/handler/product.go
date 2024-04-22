package handler

import (
	"context"

	"mxshop-go/product_svc/global"
	"mxshop-go/product_svc/model"
	"mxshop-go/product_svc/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func productModelToProtoProductInfo(product *model.Product) *proto.ProductInfo {
	return &proto.ProductInfo{
		Id:           product.ID,
		CategoryId:   product.CategoryID,
		Name:         product.Name,
		SerialNumber: product.SerialNumber,
		ClickCount:   product.ClickCount,
		SoldCount:    product.SoldCount,
		LikeCount:    product.LikeCount,
		MarketPrice:  product.MarketPrice,
		ShopPrice:    product.ShopPrice,
		Brief:        product.Brief,
		Description:  product.Description,
		FreeShipping: product.FreeShipping,
		Images:       product.Images,
		DescImages:   product.DescImages,
		FrontImage:   product.FrontImage,
		IsNew:        product.IsNew,
		IsHot:        product.IsHot,
		OnSale:       product.OnSale,
		Brand: &proto.BrandInfo{
			Id:   product.BrandID,
			Name: product.Brand.Name,
			Logo: product.Brand.Logo,
		},
		Category: &proto.CategoryBriefInfo{
			Id:   product.CategoryID,
			Name: product.Category.Name,
		},
	}
}

func (p ProductServiceServer) FilterProducts(ctx context.Context, request *proto.FilterProductsRequest) (*proto.FilterProductsResponse, error) {
	var rsp proto.FilterProductsResponse
	db := global.DB.Model(&model.Product{})

	if request.PriceMin != nil && request.GetPriceMin() > 0 {
		db = db.Where("shop_price >= ?", request.GetPriceMin())
	}
	if request.PriceMax != nil && request.GetPriceMax() > 0 {
		db = db.Where("shop_price <= ?", request.GetPriceMax())
	}
	if request.IsHot != nil {
		db = db.Where(model.Product{IsNew: request.GetIsHot()})
	}
	if request.IsNew != nil {
		db = db.Where(model.Product{IsNew: request.GetIsNew()})
	}
	if request.Brand != nil {
		db = db.Where(model.Product{BrandID: request.GetBrand()})
	}
	if request.KeyWords != nil && request.GetKeyWords() != "" {
		db = db.Where("name LIKE ?", "%"+request.GetKeyWords()+"%")
	}
	if request.Brand != nil && request.GetBrand() != 0 {
		db = db.Where(model.Product{BrandID: request.GetBrand()})
	}
	if request.TopCategory != nil && request.GetTopCategory() > 0 {
		var category model.Category
		result := global.DB.Limit(1).Find(&category, request.GetTopCategory())
		if result.Error != nil {
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		if category.Level == 1 {
			db = db.Where("category_id IN (SELECT id FROM categories WHERE upper_level_category_id in (SELECT id FROM categories WHERE upper_level_category_id = ?))", request.GetTopCategory())
		}
		if category.Level == 2 {
			db = db.Where("category_id IN (SELECT id FROM categories WHERE upper_level_category_id = ?)", request.GetTopCategory())
		}
		if category.Level == 3 {
			db = db.Where("category_id = ?", request.GetTopCategory())
		}
	}

	var products []*model.Product
	var total int64
	if result := db.Count(&total).Preload("Brand").Preload("Category").Scopes(Paginate(int(request.GetPage()), int(request.GetPageSize()))).Find(&products); result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	rsp.Total = int32(total)
	rsp.Data = make([]*proto.ProductInfo, 0, len(products))
	for _, product := range products {
		rsp.Data = append(rsp.Data, productModelToProtoProductInfo(product))
	}

	return &rsp, nil
}

func (p ProductServiceServer) BatchGetProducts(ctx context.Context, request *proto.BatchGetProductsRequest) (*proto.BatchGetProductsResponse, error) {
	var rsp proto.BatchGetProductsResponse

	if len(request.Ids) == 0 {
		rsp.Total = 0
		rsp.Data = make([]*proto.ProductInfo, 0)
		return &rsp, nil
	}

	var products []*model.Product
	var total int64
	result := global.DB.Preload("Category").Preload("Brand").Find(&products, request.Ids)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	total = result.RowsAffected

	rsp.Total = int32(total)
	rsp.Data = make([]*proto.ProductInfo, 0, len(products))
	for _, product := range products {
		rsp.Data = append(rsp.Data, productModelToProtoProductInfo(product))
	}

	return &rsp, nil
}

func (p ProductServiceServer) CreateProduct(ctx context.Context, request *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	var rsp proto.CreateProductResponse

	// Validate "Name" field
	if request.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	// Validate "CategoryID" field
	if request.CategoryId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "category_id is required")
	}
	var category model.Category
	result := global.DB.First(&category, request.CategoryId)
	if result.Error != nil {
		zap.S().Errorf("create product failed: %v", result.Error)
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "category not found")
	}

	// Validate "BrandID" field
	if request.BrandId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "brand_id is required")
	}
	var brand model.Brand
	result = global.DB.First(&brand, request.BrandId)
	if result.Error != nil {
		zap.S().Errorf("create product failed: %v", result.Error)
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "brand not found")
	}

	// Validate "CategoryID" and "BrandID" relationship
	result = global.DB.Where("category_id = ? AND brand_id = ?", request.CategoryId, request.BrandId).First(&model.CategoryBrand{})
	if result.Error != nil {
		zap.S().Errorf("create product failed: %v", result.Error)
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "category brand relationship not found")
	}

	// Validate "SerialNumber" field
	if request.SerialNumber == "" {
		return nil, status.Errorf(codes.InvalidArgument, "serial_number is required")
	}

	// Validate "MarketPrice" field
	if request.MarketPrice <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "market_price is required")
	}

	// Validate "ShopPrice" field
	if request.ShopPrice <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "shop_price is required")
	}

	// Validate "Brief" field
	if request.Brief == "" {
		return nil, status.Errorf(codes.InvalidArgument, "brief is required")
	}

	// Validate "Description" field
	if request.Description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "description is required")
	}

	// Validate "FrontImage" field
	if request.FrontImage == "" {
		return nil, status.Errorf(codes.InvalidArgument, "front_image is required")
	}

	product := model.Product{
		CategoryID:   request.CategoryId,
		BrandID:      request.BrandId,
		OnSale:       request.OnSale,
		FreeShipping: request.FreeShipping,
		IsNew:        request.IsNew,
		IsHot:        request.IsHot,
		Name:         request.Name,
		SerialNumber: request.SerialNumber,
		Brief:        request.Brief,
		Description:  request.Description,
		ClickCount:   0,
		LikeCount:    0,
		SoldCount:    0,
		MarketPrice:  request.MarketPrice,
		ShopPrice:    request.ShopPrice,
		Images:       model.StringList(request.Images),
		DescImages:   model.StringList(request.DescImages),
		FrontImage:   request.FrontImage,
	}

	if result := global.DB.Create(&product); result.Error != nil {
		zap.S().Errorf("create product failed: %v", result.Error)
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	product.Category = &category
	product.Brand = &brand
	rsp.Data = productModelToProtoProductInfo(&product)

	return &rsp, nil
}

func (p ProductServiceServer) DeleteProduct(ctx context.Context, request *proto.DeleteProductRequest) (*emptypb.Empty, error) {
	if request.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	result := global.DB.Delete(&model.Product{}, request.Id)
	if result.Error != nil {
		zap.S().Errorf("delete product failed: %v", result.Error)
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	return &emptypb.Empty{}, nil
}

func (p ProductServiceServer) UpdateProduct(ctx context.Context, request *proto.UpdateProductRequest) (*emptypb.Empty, error) {
	// TODO implement me
	return &emptypb.Empty{}, nil
}

func (p ProductServiceServer) GetProduct(ctx context.Context, request *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	var rsp proto.GetProductResponse
	var product model.Product

	if request.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	result := global.DB.Preload("Brand").Preload("Category").First(&product, request.Id)
	if result.Error != nil {
		zap.S().Errorf("get product failed: %v", result.Error)
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	rsp.Data = productModelToProtoProductInfo(&product)

	return &rsp, nil
}

package handler

import (
	"context"

	"mxshop-go/product_svc/global"
	"mxshop-go/product_svc/model"
	"mxshop-go/product_svc/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p ProductServiceServer) GetCategoryBrandList(ctx context.Context, request *proto.GetCategoryBrandListRequest) (*proto.GetCategoryBrandListResponse, error) {
	var rsp proto.GetCategoryBrandListResponse
	var categoryBrands []*model.CategoryBrand
	var total int64

	if err := global.DB.Model(&model.CategoryBrand{}).Count(&total).Error; err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	rsp.Total = int32(total)

	result := global.DB.Preload("Category").Preload("Brand").Scopes(Paginate(int(request.Page), int(request.PageSize))).Find(&categoryBrands)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	rsp.Data = make([]*proto.CategoryBrandInfo, 0, 10)
	for _, categoryBrand := range categoryBrands {
		rsp.Data = append(rsp.Data, &proto.CategoryBrandInfo{
			Id: categoryBrand.ID,
			Brand: &proto.BrandInfo{
				Id:   categoryBrand.Brand.ID,
				Name: categoryBrand.Brand.Name,
				Logo: categoryBrand.Brand.Logo,
			},
			Category: &proto.CategoryInfo{
				Id:             categoryBrand.Category.ID,
				Name:           categoryBrand.Category.Name,
				ParentCategory: *categoryBrand.Category.UpperLevelCategoryID,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.VisibleInTab,
			},
		})
	}

	return &rsp, nil
}

func (p ProductServiceServer) GetBrandsByCategory(ctx context.Context, request *proto.GetBrandsByCategoryRequest) (*proto.GetBrandsByCategoryResponse, error) {
	var rsp proto.GetBrandsByCategoryResponse
	var categoryBrands []*model.CategoryBrand

	result := global.DB.Where("category_id = ?", request.Id).Find(&categoryBrands)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	rsp.Total = int32(result.RowsAffected)
	rsp.Data = make([]*proto.BrandInfo, 0, 10)
	for _, categoryBrand := range categoryBrands {
		rsp.Data = append(rsp.Data, &proto.BrandInfo{
			Id:   categoryBrand.Brand.ID,
			Name: categoryBrand.Brand.Name,
			Logo: categoryBrand.Brand.Logo,
		})
	}

	return &rsp, nil
}

func (p ProductServiceServer) CreateCategoryBrand(ctx context.Context, request *proto.CreateCategoryBrandRequest) (*proto.CreateCategoryBrandResponse, error) {
	var rsp proto.CreateCategoryBrandResponse
	categoryBrand := model.CategoryBrand{
		CategoryID: request.CategoryId,
		BrandID:    request.BrandId,
	}

	// Check if the category and the brand are already associated
	result := global.DB.Where(map[string]interface{}{
		"category_id": request.CategoryId,
		"brand_id":    request.BrandId,
	}).First(&model.CategoryBrand{})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "the category and brand are already associated")
	}

	// Check if category exists
	result = global.DB.First(&model.Category{}, request.CategoryId)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}

	// Check if brand exists
	result = global.DB.First(&model.Brand{}, request.BrandId)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "brand not found")
	}

	// create category brand relationship
	if err := global.DB.Create(&categoryBrand).Error; err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rsp.Data = &proto.CategoryBrandInfo{
		Id: categoryBrand.ID,
	}

	return &rsp, nil
}

func (p ProductServiceServer) DeleteCategoryBrand(ctx context.Context, request *proto.DeleteCategoryBrandRequest) (*emptypb.Empty, error) {
	var categoryBrand model.CategoryBrand
	result := global.DB.Where("id = ?", request.Id).First(&categoryBrand)
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}

	if err := global.DB.Delete(&categoryBrand).Error; err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (p ProductServiceServer) UpdateCategoryBrand(ctx context.Context, request *proto.UpdateCategoryBrandRequest) (*emptypb.Empty, error) {
	var categoryBrand model.CategoryBrand

	result := global.DB.First(&categoryBrand, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category brand relationship not found")
	}

	if request.CategoryId == categoryBrand.CategoryID && request.BrandId == categoryBrand.BrandID {
		return &emptypb.Empty{}, nil
	}

	// Check if the category and the brand are already associated
	result = global.DB.Where(map[string]interface{}{
		"category_id": request.CategoryId,
		"brand_id":    request.BrandId,
	}).First(&model.CategoryBrand{})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "the category and brand are already associated")
	}

	// Check if category exists
	result = global.DB.First(&model.Category{}, request.CategoryId)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}

	// Check if brand exists
	result = global.DB.First(&model.Brand{}, request.BrandId)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "brand not found")
	}

	// Update category brand relationship
	categoryBrand.CategoryID = request.CategoryId
	categoryBrand.BrandID = request.BrandId
	if err := global.DB.Save(&categoryBrand).Error; err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

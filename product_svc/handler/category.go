package handler

import (
	"context"
	"encoding/json"

	"mxshop-go/product_svc/global"
	"mxshop-go/product_svc/model"
	"mxshop-go/product_svc/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p ProductServiceServer) GetAllCategories(_ context.Context, _ *emptypb.Empty) (*proto.GetAllCategoriesResponse, error) {
	rsp := &proto.GetAllCategoriesResponse{}
	categories := make([]*model.Category, 0)

	if err := global.DB.Where(&model.Category{Level: 1}).Preload("SubCategories.SubCategories").Find(&categories).Error; err != nil {
		return nil, err
	}

	rsp.Total = int32(len(categories))
	for _, category := range categories {
		rsp.Data = append(rsp.Data, &proto.CategoryInfo{
			Id:    category.ID,
			Name:  category.Name,
			Level: category.Level,
			IsTab: category.VisibleInTab,
		})
	}
	jsonData, err := json.Marshal(categories)
	if err != nil {
		return nil, err
	}
	rsp.JsonData = string(jsonData)

	return rsp, nil
}

func (p ProductServiceServer) GetSubCategories(_ context.Context, request *proto.GetSubCategoriesRequest) (*proto.GetSubCategoriesResponse, error) {
	rsp := proto.GetSubCategoriesResponse{}

	// get category
	var category model.Category
	result := global.DB.First(&category, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}

	// set category info
	var parentCategory int32 = 0
	if category.UpperLevelCategoryID != nil {
		parentCategory = *category.UpperLevelCategoryID
	}
	rsp.Info = &proto.CategoryInfo{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		ParentCategory: parentCategory,
		IsTab:          category.VisibleInTab,
	}

	// get sub categories
	subCategories := make([]*model.Category, 0, 10)
	if err := global.DB.Where(&model.Category{UpperLevelCategoryID: &category.ID}).Find(&subCategories).Error; err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	rsp.Total = int32(len(subCategories))
	for _, subCategory := range subCategories {
		rsp.SubCategories = append(rsp.SubCategories, &proto.CategoryInfo{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			ParentCategory: *subCategory.UpperLevelCategoryID,
			IsTab:          subCategory.VisibleInTab,
		})
	}

	return &rsp, nil
}

func (p ProductServiceServer) CreateCategory(_ context.Context, request *proto.CreateCategoryRequest) (*proto.CreateCategoryResponse, error) {
	category := model.Category{
		Name:         request.Name,
		VisibleInTab: request.IsTab,
	}
	if request.ParentCategory == 0 {
		category.Level = 1
		category.UpperLevelCategoryID = nil
	} else {
		category.UpperLevelCategoryID = &request.ParentCategory
		var upperLevelCategory model.Category
		result := global.DB.First(&upperLevelCategory, request.ParentCategory)
		if result.Error != nil {
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "upper level category not found")
		}
		category.Level = upperLevelCategory.Level + 1
	}
	if err := global.DB.Create(&category).Error; err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.CreateCategoryResponse{
		Data: &proto.CategoryInfo{
			Id:             category.ID,
			Name:           category.Name,
			ParentCategory: *category.UpperLevelCategoryID,
			Level:          category.Level,
			IsTab:          category.VisibleInTab,
		},
	}, nil
}

func (p ProductServiceServer) DeleteCategory(_ context.Context, request *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Category{}, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}

	return &emptypb.Empty{}, nil
}

func (p ProductServiceServer) UpdateCategory(_ context.Context, request *proto.UpdateCategoryRequest) (*emptypb.Empty, error) {
	var category model.Category

	// check if category exists
	result := global.DB.First(&category, request.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "category not found")
	}

	// update category
	if request.Name != "" {
		category.Name = request.Name
	}
	if request.IsTab != category.VisibleInTab {
		category.VisibleInTab = request.IsTab
	}
	if request.ParentCategory == 0 {
		category.UpperLevelCategoryID = nil
		category.Level = 1
	} else {
		category.UpperLevelCategoryID = &request.ParentCategory
		var upperLevelCategory model.Category
		result = global.DB.First(&upperLevelCategory, request.ParentCategory)
		if result.Error != nil {
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "upper level category not found")
		}
		category.Level = upperLevelCategory.Level + 1
	}
	if err := global.DB.Save(&category).Error; err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

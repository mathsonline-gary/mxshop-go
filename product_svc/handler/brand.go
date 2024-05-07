package handler

import (
	"context"

	"github.com/zycgary/mxshop-go/product_svc/global"
	"github.com/zycgary/mxshop-go/product_svc/model"
	"github.com/zycgary/mxshop-go/product_svc/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p ProductServiceServer) GetBrands(ctx context.Context, request *proto.GetBrandsRequest) (*proto.GetBrandsResponse, error) {
	rsp := proto.GetBrandsResponse{}
	brands := make([]model.Brand, 0)
	var total int64

	result := global.DB.Model(&model.Brand{}).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp.Total = int32(total)

	result = global.DB.Scopes(Paginate(int(request.Page), int(request.PageSize))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, brand := range brands {
		rsp.Data = append(rsp.Data, &proto.BrandInfo{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	return &rsp, nil
}

func (p ProductServiceServer) CreateBrand(ctx context.Context, request *proto.CreateBrandRequest) (*proto.CreateBrandResponse, error) {
	rsp := proto.CreateBrandResponse{}
	brand := model.Brand{}
	brand.Name = request.Name
	brand.Logo = request.Logo

	// Check if a brand with the same name already exists
	if result := global.DB.Where("name = ?", brand.Name).First(&model.Brand{}); result.RowsAffected > 0 {
		return nil, status.Errorf(codes.InvalidArgument, "The brand already exists")
	}

	result := global.DB.Create(&brand)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp.Data = &proto.BrandInfo{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.Logo,
	}

	return &rsp, nil
}

func (p ProductServiceServer) DeleteBrand(ctx context.Context, request *proto.DeleteBrandRequest) (*emptypb.Empty, error) {
	result := global.DB.Delete(&model.Brand{}, request.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "The brand does not exist")
	}

	return &emptypb.Empty{}, nil
}

func (p ProductServiceServer) UpdateBrand(ctx context.Context, request *proto.UpdateBrandRequest) (*emptypb.Empty, error) {
	brand := model.Brand{}

	// Check if the brand exists
	result := global.DB.First(&brand, request.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "The brand does not exist")
	}

	// Update the brand
	if request.Name != "" {
		brand.Name = request.Name
	}
	if request.Logo != "" {
		brand.Logo = request.Logo
	}
	if err := global.DB.Save(&brand).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

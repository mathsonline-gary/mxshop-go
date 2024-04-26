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

func (p ProductServiceServer) GetBanners(ctx context.Context, empty *emptypb.Empty) (*proto.GetBannersResponse, error) {
	rsp := proto.GetBannersResponse{}
	banners := make([]model.Banner, 0)

	result := global.DB.Find(&banners)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp.Total = int32(result.RowsAffected)
	for _, banner := range banners {
		rsp.Data = append(rsp.Data, &proto.BannerInfo{
			Id:    banner.ID,
			Image: banner.Image,
			Url:   banner.URL,
			Index: banner.Index,
		})
	}

	return &rsp, nil
}

func (p ProductServiceServer) CreateBanner(ctx context.Context, request *proto.CreateBannerRequest) (*proto.CreateBannerResponse, error) {
	rsp := proto.CreateBannerResponse{}
	banner := model.Banner{
		Image: request.Image,
		URL:   request.Url,
		Index: request.Index,
	}
	global.DB.Create(&banner)
	rsp.Data = &proto.BannerInfo{
		Id:    banner.ID,
		Image: banner.Image,
		Url:   banner.URL,
		Index: banner.Index,
	}

	return &rsp, nil
}

func (p ProductServiceServer) DeleteBanner(ctx context.Context, request *proto.UpdateBannerRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Banner{}, request.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "banner not found")
	}

	return &emptypb.Empty{}, nil
}

func (p ProductServiceServer) UpdateBanner(ctx context.Context, request *proto.DeleteBannerRequest) (*emptypb.Empty, error) {
	var banner model.Banner
	if result := global.DB.First(&banner, request.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "banner not found")
	}

	if request.Image != "" {
		banner.Image = request.Image
	}
	if request.Url != "" {
		banner.URL = request.Url
	}
	if request.Index != 0 {
		banner.Index = request.Index
	}

	if result := global.DB.Save(&banner); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

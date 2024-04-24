package product

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"mxshop-go/product_api/controllers"
	"mxshop-go/product_api/global"
	"mxshop-go/product_api/requests"
	"mxshop-go/product_svc/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Index(ctx *gin.Context) {
	req := proto.FilterProductsRequest{
		PriceMin:    nil,
		PriceMax:    nil,
		IsHot:       nil,
		IsNew:       nil,
		IsTab:       nil,
		TopCategory: nil,
		Page:        1,
		PageSize:    10,
		KeyWords:    nil,
		Brand:       nil,
	}

	minPrice := ctx.DefaultQuery("min_price", "0")
	minPriceInt, err := strconv.Atoi(minPrice)
	if err == nil {
		p := int32(minPriceInt)
		req.PriceMin = &p
	}

	maxPrice := ctx.DefaultQuery("max_price", "0")
	maxPriceInt, err := strconv.Atoi(maxPrice)
	if err == nil {
		p := int32(maxPriceInt)
		req.PriceMax = &p
	}

	isHot := ctx.DefaultQuery("is_hot", "-1")
	isHotInt, err := strconv.Atoi(isHot)
	if err == nil && isHotInt != -1 {
		ih := isHotInt == 1
		req.IsHot = &ih
	}

	isNew := ctx.DefaultQuery("is_new", "-1")
	isNewInt, err := strconv.Atoi(isNew)
	if err == nil && isNewInt != -1 {
		in := isNewInt == 1
		req.IsNew = &in
	}

	isTab := ctx.DefaultQuery("is_tab", "-1")
	isTabInt, err := strconv.Atoi(isTab)
	if err == nil && isTabInt != -1 {
		it := isTabInt == 1
		req.IsTab = &it
	}

	topCategory := ctx.DefaultQuery("top_category", "0")
	topCategoryInt, err := strconv.Atoi(topCategory)
	if err == nil {
		tc := int32(topCategoryInt)
		req.TopCategory = &tc
	}

	page := ctx.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err == nil && pageInt > 0 {
		req.Page = int32(pageInt)
	}

	pageSize := ctx.DefaultQuery("page_size", "10")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err == nil && pageSizeInt > 0 {
		req.PageSize = int32(pageSizeInt)
	}

	keyWords := strings.Trim(ctx.DefaultQuery("keywords", ""), " ")
	if keyWords != "" {
		req.KeyWords = &keyWords
	}

	brand := ctx.DefaultQuery("brand", "0")
	brandInt, err := strconv.Atoi(brand)
	if err == nil {
		b := int32(brandInt)
		req.Brand = &b
	}

	// Call gRPC service
	products, err := global.ProductSvcClient.FilterProducts(context.Background(), &req)
	if err != nil {
		zap.S().Errorf("grpc service FilterProducts failed: %v", err)
		controllers.HandleGRPCError(ctx, err)
		return
	}
	data := make([]interface{}, 0, len(products.Data))
	for _, product := range products.Data {
		data = append(data, map[string]interface{}{
			"id":   product.Id,
			"name": product.Name,
		})
	}
	rsp := map[string]interface{}{
		"total": products.Total,
		"data":  data,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func Store(ctx *gin.Context) {
	var req requests.StoreProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		controllers.HandleValidationError(err, ctx)
		return
	}

	rsp, err := global.ProductSvcClient.CreateProduct(context.Background(), &proto.CreateProductRequest{
		Name:         req.Name,
		SerialNumber: req.SerialNumber,
		Stocks:       req.Stock,
		MarketPrice:  req.MarketPrice,
		ShopPrice:    req.ShopPrice,
		Brief:        req.Brief,
		Description:  req.Description,
		FreeShipping: req.FreeShipping,
		Images:       req.Images,
		DescImages:   req.DescImages,
		FrontImage:   req.FrontImage,
		IsNew:        false,
		IsHot:        false,
		OnSale:       false,
		CategoryId:   req.CategoryID,
		BrandId:      req.BrandID,
	})
	if err != nil {
		controllers.HandleGRPCError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": rsp.Data.Id,
	})
}

func Show(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		return
	}

	rsp, err := global.ProductSvcClient.GetProduct(context.Background(), &proto.GetProductRequest{
		Id: int32(idInt),
	})
	if err != nil {
		controllers.HandleGRPCError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":            rsp.Data.Id,
		"name":          rsp.Data.Name,
		"brief":         rsp.Data.Brief,
		"description":   rsp.Data.Description,
		"market_price":  rsp.Data.MarketPrice,
		"shop_price":    rsp.Data.ShopPrice,
		"free_shipping": rsp.Data.FreeShipping,
		"images":        rsp.Data.Images,
		"desc_images":   rsp.Data.DescImages,
		"front_image":   rsp.Data.FrontImage,
		"category": map[string]interface{}{
			"id":   rsp.Data.Category.Id,
			"name": rsp.Data.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   rsp.Data.Brand.Id,
			"name": rsp.Data.Brand.Name,
		},
		"on_sale":     rsp.Data.OnSale,
		"is_new":      rsp.Data.IsNew,
		"is_hot":      rsp.Data.IsHot,
		"click_count": rsp.Data.ClickCount,
		"like_count":  rsp.Data.LikeCount,
		"sold_count":  rsp.Data.SoldCount,
	})
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		return
	}

	var req requests.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		controllers.HandleValidationError(err, ctx)
		return
	}

	fmt.Println(idInt, req)

	if _, err := global.ProductSvcClient.UpdateProduct(context.Background(), &proto.UpdateProductRequest{
		Id:              int32(idInt),
		Name:            req.Name,
		GoodsSn:         req.SerialNumber,
		Stocks:          req.Stock,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.Brief,
		GoodsDesc:       req.Description,
		ShipFree:        req.FreeShipping,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.FrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
		CategoryId:      req.CategoryID,
		BrandId:         req.BrandID,
	}); err != nil {
		controllers.HandleGRPCError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update product successfully",
	})
}

func Destroy(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		return
	}

	_, err = global.ProductSvcClient.DeleteProduct(context.Background(), &proto.DeleteProductRequest{
		Id: int32(idInt),
	})
	if err != nil {
		controllers.HandleGRPCError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

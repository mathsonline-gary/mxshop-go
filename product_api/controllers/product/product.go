package product

import (
	"net/http"
	"strconv"
	"strings"

	"mxshop-go/product_api/global"
	"mxshop-go/product_svc/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGRPCError(err error, ctx *gin.Context) {
	stat := http.StatusInternalServerError
	msg := "Internal error"

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				stat = http.StatusNotFound
				msg = err.Error()
			case codes.InvalidArgument:
				stat = http.StatusUnprocessableEntity
				msg = "Invalid request"
			case codes.Internal:
			default:
				stat = http.StatusInternalServerError
				msg = "Internal error"
			}
		}

		ctx.JSON(stat, gin.H{
			"message": msg,
		})
	}
}

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
	products, err := global.ProductSvcClient.FilterProducts(ctx, &req)
	if err != nil {
		zap.S().Errorf("grpc service FilterProducts failed: %v", err)
		HandleGRPCError(err, ctx)
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

package subcategory

import (
	"strconv"

	"github.com/zycgary/mxshop-go/product_api/global"
	"github.com/zycgary/mxshop-go/product_svc/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Index(ctx *gin.Context) {
	data := make([]interface{}, 0, 10)
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error("[SubCategory] [List] Invalid category ID: ", id)
		ctx.JSON(200, gin.H{
			"data": data,
			"meta": gin.H{
				"total": 0,
			},
		})
		return
	}

	rsp, err := global.ProductSvcClient.GetSubCategories(ctx, &proto.GetSubCategoriesRequest{
		Id: int32(idInt),
	})
	if err != nil {
		zap.S().Errorf("[SubCategory] [List] failed to get subcategories of category ID %d: %s", idInt, err.Error())
		ctx.JSON(200, gin.H{
			"data": data,
			"meta": gin.H{
				"total": 0,
			},
		})
		return
	}

	for _, subCategory := range rsp.SubCategories {
		data = append(data, gin.H{
			"id":             subCategory.Id,
			"name":           subCategory.Name,
			"level":          subCategory.Level,
			"parentCategory": subCategory.ParentCategory,
			"isTab":          subCategory.IsTab,
		})
	}

	ctx.JSON(200, gin.H{
		"data": data,
		"meta": gin.H{
			"total": rsp.Total,
		},
	})
}

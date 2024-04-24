package category

import (
	"context"
	"encoding/json"

	"mxshop-go/product_api/controllers"
	"mxshop-go/product_api/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Index(ctx *gin.Context) {
	rsp, err := global.ProductSvcClient.GetAllCategories(context.Background(), &emptypb.Empty{})
	if err != nil {
		controllers.HandleGRPCError(ctx, err)
	}

	data := make([]interface{}, 0, 10)
	if err := json.Unmarshal([]byte(rsp.JsonData), &data); err != nil {
		zap.S().Error("[Category] [List] failed to unmarshal json data: ", err.Error())
	}

	ctx.JSON(200, gin.H{
		"data": data,
	})
}

func Show(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"data": "show",
	})
}

package category

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"mxshop-go/product_api/controllers"
	"mxshop-go/product_api/global"
	"mxshop-go/product_api/requests"
	"mxshop-go/product_svc/proto"

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

func Create(ctx *gin.Context) {
	var req requests.StoreCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		controllers.HandleValidationError(ctx, err)
		return
	}

	rsp, err := global.ProductSvcClient.CreateCategory(context.Background(), &proto.CreateCategoryRequest{
		Name:           req.Name,
		ParentCategory: req.UpperLevelCategoryId,
		IsTab:          req.IsTab,
	})
	if err != nil {
		controllers.HandleGRPCError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": map[string]interface{}{
			"id": rsp.Data.Id,
		},
	})
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "category not found",
		})
		return
	}

	var req requests.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		controllers.HandleValidationError(ctx, err)
		return
	}

	_, err = global.ProductSvcClient.UpdateCategory(context.Background(), &proto.UpdateCategoryRequest{
		Id:             int32(idInt),
		Name:           req.Name,
		ParentCategory: req.UpperLevelCategoryId,
		IsTab:          req.IsTab,
	})
	if err != nil {
		controllers.HandleGRPCError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": map[string]interface{}{
			"id": idInt,
		},
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "category not found",
		})
		return
	}

	_, err = global.ProductSvcClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{
		Id: int32(idInt),
	})
	if err != nil {
		controllers.HandleGRPCError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

package routes

import (
	"mxshop-go/product_api/controllers/product"

	"github.com/gin-gonic/gin"
)

type ProductRouter struct {
	*gin.Engine
}

func (r *ProductRouter) Routes() {
	v1 := r.Group("/v1/products")
	v1.GET("/", product.Index)
}

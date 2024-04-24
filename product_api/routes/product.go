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

	// Get the list of products with filters
	v1.GET("/", product.Index)

	// Store a new product
	v1.POST("/", product.Store)

	// Get a product by ID
	v1.GET("/:id", product.Show)

	// Update a product by ID
	v1.PUT("/:id", product.Update)

	// Delete a product by ID
	v1.DELETE("/:id", product.Destroy)
}

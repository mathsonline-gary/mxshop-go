package routes

import (
	"github.com/zycgary/mxshop-go/product_api/controllers/product"

	"github.com/gin-gonic/gin"
)

func setProductRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/products")

	// Get the list of products with filters
	products.GET("/", product.Index)

	// Store a new product
	products.POST("/", product.Store)

	// Get a product by ID
	products.GET("/:id", product.Show)

	// Update a product by ID
	products.PUT("/:id", product.Update)

	// Delete a product by ID
	products.DELETE("/:id", product.Destroy)
}

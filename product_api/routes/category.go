package routes

import (
	"github.com/zycgary/mxshop-go/product_api/controllers/category"
	"github.com/zycgary/mxshop-go/product_api/controllers/subcategory"

	"github.com/gin-gonic/gin"
)

func setCategoryRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")

	// Get the list of categories.
	categories.GET("/", category.Index)

	// Create a new category.
	categories.POST("/", category.Create)

	// Update a category.
	categories.PUT("/:id", category.Update)

	// Delete a category.
	categories.DELETE("/:id", category.Delete)

	// Get the list of subcategories of a category.
	categories.GET("/:id/subcategories", subcategory.Index)
}

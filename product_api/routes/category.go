package routes

import (
	"mxshop-go/product_api/controllers/category"
	"mxshop-go/product_api/controllers/subcategory"

	"github.com/gin-gonic/gin"
)

func setCategoryRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")

	// Get the list of categories.
	categories.GET("/", category.Index)

	// Get the list of subcategories of a category.
	categories.GET("/:id/subcategories", subcategory.Index)
}

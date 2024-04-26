package routes

import (
	"mxshop-go/product_api/controllers/banner"

	"github.com/gin-gonic/gin"
)

func setBannerRoutes(rg *gin.RouterGroup) {
	banners := rg.Group("/banners")

	banners.GET("/", banner.Index)
	banners.POST("/", banner.Create)
	banners.PUT("/:id", banner.Update)
	banners.DELETE("/:id", banner.Delete)
}

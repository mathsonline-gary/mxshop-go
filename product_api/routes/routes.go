package routes

import "github.com/gin-gonic/gin"

func SetRoutes(r *gin.Engine) {
	rg := r.Group("/p/v1")
	setProductRoutes(rg)
	setCategoryRoutes(rg)
}

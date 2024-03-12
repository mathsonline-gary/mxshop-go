package initialize

import (
	"mxshop-go/user_api/routes"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	rg := r.Group("/v1")
	routes.InitUserRoutes(rg)
}

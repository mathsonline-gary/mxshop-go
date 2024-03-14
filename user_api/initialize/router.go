package initialize

import (
	"mxshop-go/user_api/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Router(r *gin.Engine) {
	rg := r.Group("/v1")

	zap.S().Info("initializing user routes...")
	routes.InitUserRoutes(rg)
}

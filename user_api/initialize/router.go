package initialize

import (
	"fmt"

	"mxshop-go/user_api/routes"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	fmt.Println("initializing routes...")

	rg := r.Group("/v1")
	routes.InitUserRoutes(rg)

	fmt.Println("routes initialized!")
}

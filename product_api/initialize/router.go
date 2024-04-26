package initialize

import (
	"fmt"

	"mxshop-go/product_api/routes"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	fmt.Println("routes initializing...")
	routes.SetRoutes(r)
	fmt.Println("routes initialized!")
}

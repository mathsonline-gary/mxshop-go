package initialize

import (
	"fmt"

	"mxshop-go/product_api/routes"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	fmt.Println("routes initializing...")
	router := routes.ProductRouter{Engine: r}
	router.Routes()
	fmt.Println("routes initialized!")
}

package main

import (
	"flag"
	"fmt"

	"github.com/zycgary/mxshop-go/product_api/global"
	"github.com/zycgary/mxshop-go/product_api/initialize"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ip   = flag.String("ip", "localhost", "user API IP")
	port = flag.Int("port", 8082, "user API port")
)

func main() {
	flag.Parse()

	// Init Config
	initialize.Config()
	zap.S().Debug(global.Config.ConsulConfig)
	//Init Logger
	initialize.Logger()
	// Init Router
	if global.Config.AppConfig.Env == "production" || !global.Config.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	initialize.Router(r)
	// Init gRPC client
	initialize.ProductSvcClient()

	// Start app
	fmt.Printf("starting server at %s:%d", *ip, *port)
	if err := r.Run(fmt.Sprintf("%s:%d", *ip, *port)); err != nil {
		fmt.Printf("failed to start server at %s:%d", *ip, *port)
	}
}

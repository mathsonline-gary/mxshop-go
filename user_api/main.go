package main

import (
	"flag"
	"fmt"

	"github.com/zycgary/mxshop-go/user_api/global"
	"github.com/zycgary/mxshop-go/user_api/initialize"

	"github.com/gin-gonic/gin"
)

var (
	ip   = flag.String("ip", "localhost", "user API IP")
	port = flag.Int("port", 8081, "user API port")
)

func main() {
	flag.Parse()

	// Init Config
	initialize.Config()

	// Init Logger
	initialize.Logger()

	// Init router
	if global.Config.AppConfig.Env == "production" || !global.Config.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	initialize.Router(r)

	// Init user v1 service
	initialize.InitUserSvcClient()

	// Start app
	fmt.Printf("starting server at %s:%d", *ip, *port)
	if err := r.Run(fmt.Sprintf("%s:%d", *ip, *port)); err != nil {
		fmt.Printf("failed to start server at %s:%d", *ip, *port)
	}

}

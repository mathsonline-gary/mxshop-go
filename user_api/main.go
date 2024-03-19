package main

import (
	"flag"
	"fmt"

	"mxshop-go/user_api/initialize"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ip   = flag.String("ip", "localhost", "user API IP")
	port = flag.Int("port", 8081, "user API port")
)

func main() {
	flag.Parse()

	// Init Logger
	initialize.Logger()

	// Init Config
	initialize.Config()

	// Init router
	r := gin.Default()
	initialize.Router(r)

	// Init user grpc service
	initialize.InitUserSvcClient()

	// Start app
	zap.S().Debugf("starting server at %s:%d", *ip, *port)
	if err := r.Run(fmt.Sprintf("%s:%d", *ip, *port)); err != nil {
		zap.S().Panicf("failed to start server at %s:%d", *ip, *port)
	}

}

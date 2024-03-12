package main

import (
	"flag"
	"fmt"

	"mxshop-go/user_api/initialize"

	"github.com/gin-gonic/gin"
)

var (
	ip   = flag.String("ip", "localhost", "user API IP")
	port = flag.Int("port", 8081, "user API port")
)

func main() {
	flag.Parse()

	// Init Logger
	initialize.Logger()

	// Init router
	r := gin.Default()
	initialize.Router(r)

	// Start app
	if err := r.Run(fmt.Sprintf("%s:%d", *ip, *port)); err != nil {
		fmt.Printf("Failed to start server at %s:%d\r\n", *ip, *port)
	}

}

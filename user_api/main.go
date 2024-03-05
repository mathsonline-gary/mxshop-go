package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Start app
	port := 8088
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Failed to start server at port: %d", port)
	}

}

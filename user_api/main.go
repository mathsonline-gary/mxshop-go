package main

import (
	"fmt"

	"mxshop-go/user_api/client_stub"

	"github.com/gin-gonic/gin"
)

func main() {
	client := client_stub.NewHelloServiceClient("tcp", "localhost:9090")

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		var reply string
		_ = client.Hello("Gary", &reply)

		c.JSON(200, gin.H{
			"message": reply,
		})
	})

	// Start app
	port := 8080
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Failed to start server at port: %d", port)
	}

}

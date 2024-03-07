package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "user_api",
		})
	})

	// Start app
	port := 8081
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Failed to start server at port: %d\r\n", port)
	}

}

package routes

import (
	"github.com/zycgary/mxshop-go/user_api/controllers"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("/", controllers.Index)
}

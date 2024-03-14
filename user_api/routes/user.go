package routes

import (
	"mxshop-go/user_api/http/controllers/user_controller"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("/", user_controller.Index)
}

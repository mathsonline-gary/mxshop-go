package routes

import (
	"mxshop-go/user_api/http/controllers/user_controller"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(rg *gin.RouterGroup) {
	urg := rg.Group("/users")
	urg.GET("", user_controller.Index)
}

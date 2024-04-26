package banner

import "github.com/gin-gonic/gin"

func Index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Banner Index",
	})
}

func Create(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Banner Create",
	})
}

func Update(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Banner Update",
	})
}

func Delete(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Banner Delete",
	})
}

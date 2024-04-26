package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGRPCError(ctx *gin.Context, err error) {
	stat := http.StatusInternalServerError
	msg := "Internal error"

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				stat = http.StatusNotFound
				msg = e.Message()
			case codes.InvalidArgument:
				stat = http.StatusUnprocessableEntity
				msg = "Invalid request"
				if e.Message() != "" {
					msg += ": " + e.Message()
				}
			case codes.Internal:
			default:
				stat = http.StatusInternalServerError
				msg = "Internal error"
			}
		}

		ctx.JSON(stat, gin.H{
			"message": msg,
		})
	}
}

func HandleValidationError(ctx *gin.Context, err error) {
	zap.S().Error("Validation error: ", err)
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": err.Error(),
	})
}

package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWT(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get tokenString from header.
		tokenString := extractTokenString(ctx)
		if tokenString == "" {
			unauthorized(ctx)
			return
		}

		// Validate tokenString.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			unauthorized(ctx)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauthorized(ctx)
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			unauthorized(ctx)
			return
		}
		if float64(time.Now().Unix()) > exp {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token expired"})
			return
		}

		// Set id to context.
		idFloat, ok := claims["id"].(float64)
		if !ok {
			unauthorized(ctx)
			return
		}
		ctx.Set("id", int(idFloat))

		ctx.Next()
	}
}

func extractTokenString(ctx *gin.Context) string {
	token := ctx.Query("token")
	if token != "" {
		return token
	}

	bearerToken := ctx.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func unauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
}

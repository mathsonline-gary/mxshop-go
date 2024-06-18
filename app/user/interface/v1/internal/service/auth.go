package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/logic"
	"github.com/zycgary/mxshop-go/pkg/log"
)

type AuthService struct {
	auc    *logic.AuthUseCase
	logger *log.Sugar
}

func NewAuthService(auc *logic.AuthUseCase, logger log.Logger) *AuthService {
	return &AuthService{
		auc:    auc,
		logger: log.NewSugar(logger),
	}
}

func (s *AuthService) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request.",
		})
		return
	}

	// Validate email: required, email format, length 3-255
	if req.Email == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Email is required.",
		})
		return
	}

	// TODO: Validate email format

	if len(req.Email) < 3 || len(req.Email) > 255 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Email length must be between 3 and 255.",
		})
		return
	}

	// Validate password: required, length 6-255
	if req.Password == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Password is required.",
		})
		return
	}
	if len(req.Password) < 6 || len(req.Password) > 255 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Password length must be between 6 and 255.",
		})
		return
	}

	token, err := s.auc.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Login failed.",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"error":   false,
		"message": "Login success.",
		"data": map[string]any{
			"token": token,
		},
	})
}

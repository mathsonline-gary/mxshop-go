package service

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/logic"
	"github.com/zycgary/mxshop-go/pkg/log"
)

type AuthService struct {
	auc       *logic.AuthUseCase
	validator *validator.Validate
	logger    *log.Sugar
}

func NewAuthService(auc *logic.AuthUseCase, logger log.Logger) *AuthService {
	return &AuthService{
		auc:       auc,
		validator: validator.New(),
		logger:    log.NewSugar(logger),
	}
}

func (s *AuthService) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email,omitempty" validate:"required,email,min=3,max=255"`
		Password string `json:"password,omitempty" validate:"required,min=6,max=255"`
	}

	// Bind request.
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request.",
		})
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Validate request.
	if err := s.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request.",
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

package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/logic"
	"github.com/zycgary/mxshop-go/pkg/log"
)

type UserService struct {
	uuc    *logic.UserUseCase
	logger *log.Sugar
}

func NewUserService(uuc *logic.UserUseCase, logger log.Logger) *UserService {
	return &UserService{
		uuc:    uuc,
		logger: log.NewSugar(logger),
	}
}

func (s *UserService) Index(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")

	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	ps, err := strconv.Atoi(pageSize)
	if err != nil {
		ps = 10
	}

	// Call use case.
	ul, err := s.uuc.GetList(ctx, int32(p), int32(ps))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error":   true,
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"error":   false,
		"message": "success",
		"data": map[string]any{
			"total": ul.Total,
			"data":  ul.Data,
		},
	})
}

package logic

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zycgary/mxshop-go/pkg/log"
)

type AuthUseCase struct {
	ur     UserRepository
	secret string
	logger *log.Sugar
}

func NewAuthUseCase(secret string, ur UserRepository, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		ur:     ur,
		secret: secret,
		logger: log.NewSugar(logger),
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil || user.ID <= 0 {
		uc.logger.Error("[AuthUseCase] [Login]: User not found")
		return "", errors.New("user not found")
	}

	ok, err := uc.ur.CheckPassword(ctx, password, user.Password)
	if err != nil {
		return "", err
	}
	if !ok {
		uc.logger.Error("[AuthUseCase] [Login]: Password mismatch")
		return "", errors.New("password mismatch")
	}

	// Generate JWT token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(uc.secret))
	if err != nil {
		uc.logger.Errorf("[AuthUseCase] [Login]: Generate JWT token. %v", err)
		return "", err
	}

	return tokenString, nil
}

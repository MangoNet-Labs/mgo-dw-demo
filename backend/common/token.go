package common

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
	"user/model"
)

func GenerateToken(user model.Address, secret string, expireSeconds int64) (string, error) {
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"userId":     user.ID,
		"mgoAddress": user.MgoAddress,
		"solAddress": user.SolanaAddress,
		"exp":        now + expireSeconds,
		"iat":        now,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

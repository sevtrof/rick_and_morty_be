package service

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	secretKey string
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{secretKey: secretKey}
}

func (s *TokenService) GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   strconv.Itoa(userID),
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

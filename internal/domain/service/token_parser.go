package service

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

type JwtTokenParser struct {
	secret string
}

func NewJwtTokenParser(secret string) *JwtTokenParser {
	return &JwtTokenParser{
		secret: secret,
	}
}

func (p *JwtTokenParser) ParseToken(token string) (claims *jwt.StandardClaims, err error) {
	token = strings.TrimSpace(token)

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

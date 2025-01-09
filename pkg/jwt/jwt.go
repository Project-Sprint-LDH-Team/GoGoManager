package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Maker interface {
	GenerateToken(userID uint) (string, error)
	VerifyToken(token string) (*Claims, error)
}

type JWTMaker struct {
	secretKey string
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTMaker(secretKey string) Maker {
	return &JWTMaker{secretKey: secretKey}
}

func (maker *JWTMaker) GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(maker.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

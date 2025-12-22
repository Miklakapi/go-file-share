package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalid  = errors.New("token invalid")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenParse    = errors.New("token parse error")
	ErrTokenSignAlgo = errors.New("unexpected signing method")
)

type Claims struct {
	jwt.RegisteredClaims
}

func GenerateJWT(expMinutes int16, secret []byte) (string, error) {
	now := time.Now()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(expMinutes))),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateJWT(tokenString string, secret []byte) error {
	claims := &Claims{}

	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSignAlgo
		}
		return secret, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return ErrTokenExpired
		default:
			return fmt.Errorf("%w: %v", ErrTokenParse, err)
		}
	}

	if parsed == nil || !parsed.Valid {
		return ErrTokenInvalid
	}

	return nil
}

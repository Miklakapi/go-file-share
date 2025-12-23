package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secret []byte
	now    func() time.Time
}

var _ ports.TokenService = (*JwtService)(nil)

func NewJwtService(secret []byte) *JwtService {
	return &JwtService{secret: secret, now: time.Now}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (s *JwtService) Issue(ctx context.Context, ttl time.Duration) (string, time.Time, error) {
	if err := ctx.Err(); err != nil {
		return "", time.Time{}, err
	}

	now := s.now()
	exp := now.Add(ttl)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(s.secret)
	return token, exp, err
}

func (s *JwtService) Validate(ctx context.Context, tokenString string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ports.ErrTokenSignAlgo
		}
		return s.secret, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return ports.ErrTokenExpired
		default:
			return fmt.Errorf("%w: %v", ports.ErrTokenParse, err)
		}
	}

	if parsed == nil || !parsed.Valid {
		return ports.ErrInvalidToken
	}

	return nil
}

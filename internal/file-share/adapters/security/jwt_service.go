package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtService struct {
	secret []byte
	now    func() time.Time
}

func NewJwtService(secret []byte) *JwtService {
	return &JwtService{secret: secret, now: time.Now}
}

type Claims struct {
	RoomID string
	jwt.RegisteredClaims
}

func (s *JwtService) Issue(ctx context.Context, roomID uuid.UUID, ttl time.Duration) (string, time.Time, error) {
	if err := ctx.Err(); err != nil {
		return "", time.Time{}, err
	}

	now := s.now()
	exp := now.Add(ttl)

	claims := Claims{
		RoomID: roomID.String(),
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
	_, err := s.parseClaims(ctx, tokenString)
	return err
}

func (s *JwtService) ValidateWithRoom(ctx context.Context, roomID uuid.UUID, tokenString string) error {
	claims, err := s.parseClaims(ctx, tokenString)
	if err != nil {
		return err
	}

	if claims.RoomID != "" && claims.RoomID != roomID.String() {
		return ports.ErrTokenRoomMismatch
	}

	return nil
}

func (s *JwtService) parseClaims(ctx context.Context, tokenString string) (*Claims, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
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
			return nil, ports.ErrTokenExpired
		default:
			return nil, fmt.Errorf("%w: %v", ports.ErrTokenParse, err)
		}
	}

	if parsed == nil || !parsed.Valid {
		return nil, ports.ErrInvalidToken
	}

	return claims, nil
}

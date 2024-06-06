package transport

import (
	"birthdayNotification/internal/models"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	signingKey = "safasGdf123fgdfg1SFa"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

func (t *Transport) CreateUser(ctx context.Context, u *models.User) error {
	return t.storage.NewUser(ctx, u)
}

func (t *Transport) GenerateToken(ctx context.Context, u *models.User) (string, error) {
	id, err := t.storage.Auth(ctx, u)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})

	return token.SignedString([]byte(signingKey))
}

func (t *Transport) ParseToken(ctx context.Context, accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims.UserId, nil
	} else {
		return "", errors.New("token claims are not type *tokenClaims or not Valid")
	}
}

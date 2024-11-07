package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var ErrInvalidToken = errors.New("invalid token")

type claims struct {
	jwt.RegisteredClaims
	UserGUID string `json:"guid"`
}

type JWTBuilder struct {
	secretKey string
	tokenExp  time.Duration
}

func NewJWTBuilder(secretKey string, tokenExp time.Duration) *JWTBuilder {
	return &JWTBuilder{
		secretKey: secretKey,
		tokenExp:  tokenExp,
	}
}

func (b JWTBuilder) BuildJWTString(userGUID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(b.tokenExp)),
		},

		UserGUID: userGUID,
	})

	return token.SignedString([]byte(b.secretKey))
}

func (b JWTBuilder) GetUserGUID(tokenString string) (string, error) {
	c := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(b.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	return c.UserGUID, nil
}

package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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

	tokenString, err := token.SignedString([]byte(b.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

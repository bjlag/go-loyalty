package auth_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
)

func TestJWTBuilder_BuildJWTString(t *testing.T) {
	secretKey := "secret"
	userGUID := "41d2f86c-6ce5-4732-a485-6d09d7a9b3f7"

	t.Run("success", func(t *testing.T) {
		b := auth.NewJWTBuilder(secretKey, time.Hour)
		got, err := b.BuildJWTString(userGUID)
		require.NoError(t, err)

		token, err := jwt.Parse(got, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secretKey), nil
		})
		require.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)
		assert.Equal(t, userGUID, claims["guid"])
		assert.True(t, claims.VerifyExpiresAt(time.Now().Unix(), true))
	})
}

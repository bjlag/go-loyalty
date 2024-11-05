package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
)

func TestHasher_HashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		password := "123456"
		h := auth.NewHasher()
		got, err := h.HashPassword(password)

		assert.NoError(t, err)
		assert.True(t, h.ComparePasswords(got, password))
	})

	t.Run("failed", func(t *testing.T) {
		h := auth.NewHasher()
		got, err := h.HashPassword("123456")

		require.NoError(t, err)
		assert.False(t, h.ComparePasswords(got, "654321"))
	})
}

package auth

import (
	"context"
	"errors"
)

var ErrUserGUIDNotFound = errors.New("user GUID not found in context")

func UserGUIDFromContext(ctx context.Context) (string, error) {
	switch v := ctx.Value("userGUID").(type) {
	case string:
		return v, nil
	}

	return "", ErrUserGUIDNotFound
}

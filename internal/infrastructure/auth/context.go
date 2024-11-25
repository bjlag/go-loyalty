package auth

import (
	"context"
	"errors"
)

type ctxKeyUserGUID int

const UserGUIDKey ctxKeyUserGUID = 0

var ErrUserGUIDNotFound = errors.New("user GUID not found in context")

func UserGUIDFromContext(ctx context.Context) (string, error) {
	switch v := ctx.Value(UserGUIDKey).(type) {
	case string:
		return v, nil
	}

	return "", ErrUserGUIDNotFound
}

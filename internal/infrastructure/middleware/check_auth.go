package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
)

func CheckAuth(jwt *auth.JWTBuilder, log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			token := strings.Replace(authHeader, "Bearer ", "", 1)

			userGUID, err := jwt.GetUserGUID(token)
			if err != nil {
				if !errors.Is(err, auth.ErrInvalidToken) {
					log.WithError(err).Error("Failed to validate token")
				}

				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userGUID", userGUID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
